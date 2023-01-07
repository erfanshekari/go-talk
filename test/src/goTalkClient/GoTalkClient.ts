import GoTalkBase from "./GoTalkBase"
import Client from "./Client"


class GoTalkClient extends GoTalkBase implements Client {
    async connect() {}
    async reconnect() {}
    async close() {}
}

export default GoTalkClient