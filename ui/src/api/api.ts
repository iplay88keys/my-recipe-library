import axios from "axios";

export default class Api {
    static async get(url: string, auth: boolean = true) {
        let token = localStorage.getItem("access_token") || null;

        let config = {
            headers: {
                "Accept": "application/json",
                ...(auth && token && {"Authorization": "bearer " + token})
            }
        };

        return await axios.create({
            validateStatus: function (status) {
                return status === 200;
            }
        }).get(url, config);
    }

    static async post(url: string, body: string, auth: boolean = true) {
        let token = localStorage.getItem("access_token") || null;

        let config = {
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json",
                ...(auth && token && {"Authorization": "bearer " + token})
            }
        };

        return await axios.create({
            validateStatus: function (status) {
                return status === 200;
            }
        }).post(url, body, config);
    }
}
