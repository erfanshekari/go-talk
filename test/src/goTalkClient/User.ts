import jwt_decode from "jwt-decode";

type JWTAccessTokenPayload = {
    exp : number
    iat : number
    jti : string
    token_type : string
    user_id : number
}

class User {
    ID:string
    constructor(accessToken:string) {
        let decoded: JWTAccessTokenPayload = jwt_decode(accessToken)
        this.ID = String(decoded.user_id)
    }
}

export default User