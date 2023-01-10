import GoTalkTypes from "./types";
import Encrypt from "./Encrypt";
import User from "./User";

class GoTalkBase {
  private encrypt: Encrypt; // ill test this later
  private state: GoTalkTypes.State;
  private config: GoTalkTypes.Config;
  private socket?: WebSocket;
  private user?: User;
  constructor(config: GoTalkTypes.Config) {
    this.config = config;
    this.state = {
      initializing: true,
      closed: false,
      connected: false,
      connecting: false,
      initialized: false,
      isKeyExchangeDone: false,
      authenticated: false,
    };
    this.initUser();
    let instance = this;
    this.encrypt = new Encrypt({
      config: this.config,
      onSuccess: () => {
        instance.setState((state) => ({ ...state, initialized: true, initializing: false }));
      },
      isExchangeDone: () => instance.state.isKeyExchangeDone,
    });
  }

  private async initUser() {
    this.user = new User(await this.config.accessToken());
  }

  setState(f: GoTalkTypes.SetState) {
    this.state = f(this.state);
    if (this.config.onStateChange) {
      this.config.onStateChange(this.state);
    }
  }

  setOnStateChangeFunction(f: (state: GoTalkTypes.State) => void) {
    this.config.onStateChange = f;
    f(this.state);
  }

  async connect() {
    if (!this.state.initialized || this.state.initializing) {
      throw Error(
        "You can't connect to server before GoTalkClient initialization."
      );
    }
    this.setState((state) => ({ ...state, connecting: true }));
    this.socket = new WebSocket(this.config.ws);
    this.registerListeners();
  }

  async close() {
    this.setState((state) => ({ ...state, connected: false, closed: true }));
    this.socket?.close();
    this.removeListeners();
  }

  private registerListeners() {
    if (this.socket) {
      this.socket.addEventListener<"open">("open", (e) =>
        this.onOpenHandler(this, e)
      );
      this.socket.addEventListener<"message">("message", (e) =>
        this.onMessageHandler(this, e)
      );
      this.socket.addEventListener<"error">("error", (e) =>
        this.onErrorHandler(this, e)
      );
      this.socket.addEventListener<"close">("close", (e) =>
        this.onCloseHandler(this, e)
      );
    }
  }

  private removeListeners() {
    if (this.socket) {
      this.socket.removeEventListener<"open">("open", (e) =>
        this.onOpenHandler(this, e)
      );
      this.socket.removeEventListener<"message">("message", (e) =>
        this.onMessageHandler(this, e)
      );
      this.socket.removeEventListener<"error">("error", (e) =>
        this.onErrorHandler(this, e)
      );
      this.socket.removeEventListener<"close">("close", (e) =>
        this.onCloseHandler(this, e)
      );
    }
  }

  private async authenticate() {
    let authToken = {
      accessToken: await this.config.accessToken(),
    };
    let authJsonBinaryPayload = this.encrypt.encrypt(authToken);
    if (authJsonBinaryPayload) {
      let l = JSON.stringify(authJsonBinaryPayload);
      console.log(l);
      this.socket?.send(l);
    }
  }

  private onOpenHandler(instance: GoTalkBase, event: Event) {
    let clientPublicKey: GoTalkTypes.Events.ClientPublicKey = {
        publicKey: instance.encrypt.getPublicKey(),
      }
    instance.socket?.send(JSON.stringify(clientPublicKey));
    this.setState((state) => ({
      ...state,
      connected: true,
      connecting: false,
    }));
  }

  private onMessageHandler(instance: GoTalkBase, event: MessageEvent<any>) {
    console.log(event)
    const e: GoTalkTypes.Message = JSON.parse(event.data);
    if (!instance.encrypt.isExchangeDone()) {
      let serverPublicKey: GoTalkTypes.Events.ServerPublicKey =
        instance.encrypt.decrypt(e);

      if (serverPublicKey.publicKey) {
        instance.encrypt.setServerPublicKey(serverPublicKey.publicKey);

        instance.setState((state) => ({ ...state, isKeyExchangeDone: true }));

        instance.authenticate();
        return;
      } else {
        throw Error("RSA Key Exchange Faield...");
      }
    } else if (!instance.state.authenticated) {
      let userACK: GoTalkTypes.Events.ServerUserACK =
        instance.encrypt.decrypt(e);
      if (userACK && userACK.userID && userACK.userID === instance.user?.ID) {
        instance.setState((state) => ({ ...state, authenticated: true }));
        console.log("user is now authenticated", instance.user);
      }
    }
  }

  private onErrorHandler(instance: GoTalkBase, event: Event) {
    console.log(event, typeof event);
  }

  private onCloseHandler(instance: GoTalkBase, event: CloseEvent) {
    console.log(event, typeof event);
    this.setState((state) => ({
      ...state,
      connecting: false,
      closed: true,
      connected: false,
    }));
    instance.close();
  }
}

export default GoTalkBase;
