import { default as GTCTypes } from "./types"
import JSEncrypt from "jsencrypt";


const keySize = 501

interface RSA {
    client: ClientNode
    server?: ServerNode
    setServer: (publicKey: string) => void
    exchanged: () => boolean
    dropServer: () => void
}

interface Server {
    encrypt: (msg: object) => GTCTypes.Message | boolean
}

class ServerNode implements Server {

    ecyt: JSEncrypt

    constructor(publicKey: string) {
        this.ecyt = new JSEncrypt()
        this.ecyt.setPublicKey(publicKey)
    }

    encrypt(msg: object) {
        let msgString = JSON.stringify(msg)
        const msgLen = msgString.length
        if (msgLen > keySize) {
            var encryptedLen = 0
            var encryptedArray: string[] = []
            for (let i = 0; i < Math.ceil(msgLen / keySize); i++) {
                if ((encryptedLen + keySize) > msgLen) {
                    // last part smaller then max size
                    let target = msgString.slice(encryptedLen, (msgLen - encryptedLen))
                    let encrypted = this.ecyt.encrypt(target)
                    if (encrypted) {
                        encryptedArray = [...encryptedArray, encrypted]
                    } else {
                        return false
                    }
                    break
                } else if (i === 0) {
                    // first part
                    let target = msgString.slice(0, keySize)
                    let encrypted = this.ecyt.encrypt(target)
                    if (encrypted) {
                        encryptedArray = [...encryptedArray, encrypted]
                    } else {
                        return false
                    }
                    encryptedLen += keySize
                    continue
                }
                let target = msgString.slice(encryptedLen, (encryptedLen + keySize))
                let encrypted = this.ecyt.encrypt(target)
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
            } as GTCTypes.Message
        } else {
            let encrypted = this.ecyt.encrypt(msgString)
            if (encrypted) {
                return {
                    type: "byte",
                    content: encrypted
                } as GTCTypes.Message
            }
        }
        return false
    }

}

interface Client {
    decrypt: (msg: GTCTypes.Message) => object
}

class ClientNode implements Client {

    privateKey: GTCTypes.PrivateKey

    encrypt: JSEncrypt

    constructor(privateKey: GTCTypes.PrivateKey) {
        this.privateKey = privateKey
        this.encrypt = new JSEncrypt()
        this.encrypt.setPublicKey(privateKey.publicKey)
        this.encrypt.setPrivateKey(privateKey.privateKey)
    }

    decrypt(msg: GTCTypes.Message) {
        var decrypted: string = ""
        switch(msg.type) {
            case GTCTypes.MessageType.byte:
                let response = this.encrypt.decrypt(msg.content as string)
                if (response) {
                    decrypted = response
                } else {
                    throw Error("Decryption Failed...")
                }
                break
            case GTCTypes.MessageType.byteArray:
                for (const val of msg.content) {
                    let response = this.encrypt.decrypt(val)
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

class RSANode implements RSA {
    private _: JSEncrypt = new JSEncrypt() // init with no use for first time, recommended by jsEncryot developers

    client: ClientNode
    server?: ServerNode

    constructor(privateKey: GTCTypes.PrivateKey) {
        this.client = new ClientNode(privateKey)
    }

    setServer(publicKey: string) {
        this.server = new ServerNode(publicKey)
    }

    exchanged() {
        return Boolean(this.client) && Boolean(this.server)
    }

    dropServer() {
        this.server = undefined
    }
}

type CreateAsyncRSA = (config: GTCTypes.Config) => Promise<RSANode>

export const createAsyncRSA: CreateAsyncRSA = async (config: GTCTypes.Config) => {
    var privateKey: GTCTypes.PrivateKey
    const controller = new AbortController();
            const timeOut = setTimeout(() => controller.abort(), 3000);
            const response = await fetch(
                config.rest + "/genkey", 
                {
                    method: "POST", 
                    headers: {
                        "Authorization": ("Bearer " + (await config.accessToken()))
                    },
                    
            })
            .catch(e => e.response)
            clearTimeout(timeOut)
            if (response.status === 200) {
                privateKey = await response.json()
                return new RSANode(privateKey)
            } else {
                throw Error(`Can't get keys from api`)
            }
}

type CreateRSA = (privateKey: GTCTypes.PrivateKey) => RSA

export const createRSA: CreateRSA = (privateKey: GTCTypes.PrivateKey) => {
    return new RSANode(privateKey)
}

export default RSANode