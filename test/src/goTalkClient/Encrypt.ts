import JSEncrypt from "jsencrypt";
import GoTalkTypes from "./types"


const keySize = 501

type EncryptInput = {
    config: GoTalkTypes.Config
    onSuccess: () => void
}

class EncryptBase {

    jsEncrypt: JSEncrypt = new JSEncrypt({
        log: true
    })

    privateKey?: GoTalkTypes.PrivateKey
    
    private OnSuccessFunc: () => void

    serverPublicKey?: string

    constructor(input:EncryptInput) {
        this.OnSuccessFunc = input.onSuccess
        if (input.config.privateKey) {
            this.setKeysFromConfig(input.config)
        } else {
            this.genNewKeys(input.config)
        }
    }


    private setKeysFromConfig(conf: GoTalkTypes.Config) {
        if (conf.privateKey) {
            this.privateKey = conf.privateKey
            this.jsEncrypt.setPublicKey(conf.privateKey?.publicKey!)
            this.jsEncrypt.setPrivateKey(conf.privateKey?.privateKey!)
            this.OnSuccessFunc()
        } else throw Error("Private Key is not in input config")
    }

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
            this.jsEncrypt.setPublicKey(this.privateKey?.publicKey!)
            this.jsEncrypt.setPrivateKey(this.privateKey?.privateKey!)
            this.OnSuccessFunc()
    }

}

interface EncryptInterface {
    encrypt: (msg: object) => GoTalkTypes.Message | boolean
    decrypt: (msg: GoTalkTypes.Message) => object
    setServerPublicKey: (key: string) => void
    getPublicKey: () => string
    exchangeIsDone: () => boolean
}

class Encrypt extends EncryptBase implements EncryptInterface {
    setServerPublicKey(key: string) {
        this.serverPublicKey = key
    }

    getPublicKey() {
        return this.privateKey?.publicKey || ""
    }

    exchangeIsDone() {
        return Boolean(this.serverPublicKey)
    }

    encrypt(msg: object) {
        if (!this.exchangeIsDone()) throw Error("You Can't encrypt before key exchange")
        let msgString = JSON.stringify(msg)
        const encrypt = new JSEncrypt()
        encrypt.setPublicKey(this.serverPublicKey!)
        const msgLen = msgString.length
        if (msgLen > keySize) {
            var encryptedLen = 0
            var encryptedArray: string[] = []
            for (let i = 0; i < Math.ceil(msgLen / keySize); i++) {
                if ((encryptedLen + keySize) > msgLen) {
                    // last part smaller then max size
                    let target = msgString.slice(encryptedLen, (msgLen - encryptedLen))
                    let encrypted = encrypt.encrypt(target)
                    if (encrypted) {
                        encryptedArray = [...encryptedArray, encrypted]
                    } else {
                        return false
                    }
                    break
                } else if (i === 0) {
                    // first part
                    let target = msgString.slice(0, keySize)
                    let encrypted = encrypt.encrypt(target)
                    if (encrypted) {
                        encryptedArray = [...encryptedArray, encrypted]
                    } else {
                        return false
                    }
                    encryptedLen += keySize
                    continue
                }
                let target = msgString.slice(encryptedLen, (encryptedLen + keySize))
                let encrypted = encrypt.encrypt(target)
                if (encrypted) {
                    encryptedArray = [...encryptedArray, encrypted]
                } else {
                    return false
                }
                encryptedLen += keySize
            }
            return {
                type: "byteArray",
                content: encryptedArray
            } as GoTalkTypes.Message
        } else {
            let encrypted = encrypt.encrypt(msgString)
            if (encrypted) {
                return {
                    type: "byte",
                    content: encrypted
                } as GoTalkTypes.Message
            }
        }
        return false
    }

    decrypt(msg: GoTalkTypes.Message) {
        var decrypted: string = ""
        switch(msg.type) {
            case GoTalkTypes.MessageType.byte:
                let response = this.jsEncrypt.decrypt(msg.content as string)
                if (response) {
                    decrypted = response
                } else {
                    throw Error("Decryption Failed...")
                }
                break
            case GoTalkTypes.MessageType.byteArray:
                for (const val of msg.content) {
                    let response = this.jsEncrypt.decrypt(val)
                    if (response) {
                        decrypted += response
                    } else {
                        throw Error("Decryption Failed...")
                    }
                }
                break
        }
        return JSON.parse(decrypted)
    }
}

export default Encrypt
