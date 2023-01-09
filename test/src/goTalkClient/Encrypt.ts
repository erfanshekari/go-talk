import JSEncrypt from "jsencrypt";
import GoTalkTypes from "./types"


class Encrypt {
    private encrypt: JSEncrypt = new JSEncrypt({
        log: true
    })
    private privateKey?: GoTalkTypes.PrivateKey

    private async genNewKeys(conf: GoTalkTypes.Config) {
        const controller = new AbortController();
            const timeOut = setTimeout(() => controller.abort(), 3000);
            const response = await fetch(
                conf.rest + "/genkey", 
                {
                    method: "POST", 
                    headers: {
                        "Authorization": (await conf.accessToken())
                    },
                    
            })
            .catch(e => e.response)
            clearTimeout(timeOut)
            if (response.status === 200) {
                this.privateKey = await response.json()
            } else {
                throw Error(`Can't get keys from api`)
            }
            this.encrypt.setPublicKey(this.privateKey?.publicKey!)
            this.encrypt.setPrivateKey(this.privateKey?.privateKey!)
        }
}


export default Encrypt
