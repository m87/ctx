import axios from "axios";
import {ContextApi} from "@/api/api-context";

export const httpConfig = {
    baseURL: "/api",
    withCredentials: true,
    withXSRFToken: true,
    timeout: 6000,
    headers: {
        Accept: "application/json",
    },
};

export const http = axios.create(httpConfig);

export class Api {
    context = new ContextApi();
}

export const api = new Api();
