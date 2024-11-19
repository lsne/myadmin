import { http } from "@/utils/http";
import { baseUrlApi } from "./utils";

type Result = {
  code: number;
  message: string;
  data: Array<any>;
};

/** 后端请求左侧菜单栏 */
export const getAsyncRoutes = () => {
  return http.request<Result>("post", baseUrlApi("route/getAsyncRoutes"));
};
