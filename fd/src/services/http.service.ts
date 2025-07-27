/**
 * 封装HTTP请求服务
 * @param url 请求URL
 * @param method 请求方法
 * @param data 请求数据
 * @returns 响应数据
 */
// 基础URL配置
const BASE_URL = 'http://localhost:8080';

export async function request(
    endpoint: string,
    method: string = "GET",
    data: any = null
): Promise<string> {
    const url = `${BASE_URL}${endpoint}`;
    const options: RequestInit = {
        method,
        headers: {
            "Content-Type": "application/json"
        }
    };

    if (data) {
        options.body = JSON.stringify(data);
    }

    try {
        const response = await fetch(url, options);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.text();
    } catch (error) {
        console.error("Request failed:", error);
        throw error;
    }
}
