import GoTalkBase from "./GoTalkBase"
import Client from "./Client"


class GoTalkClient extends GoTalkBase implements Client {
    async connect() {
        await super.connect()
    }
    async reconnect() {
    }
    async close() {
        await super.close()
    }
}

export default GoTalkClient