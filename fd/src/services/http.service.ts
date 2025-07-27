/**
 * 封装HTTP请求服务
 * @param url 请求URL
 * @param method 请求方法
 * @param data 请求数据
 * @returns 响应数据
 */
export async function request(
    url: string,
    method: string = "GET",
    data: any = null
) {
    const options: RequestInit = {
        method,
        headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
        },
        credentials: "include", // 包含cookies
    };

    if (data) {
        options.body = JSON.stringify(data);
    }

    try {
        const response = await fetch(url, options);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.json();
    } catch (error) {
        console.error("Request failed:", error);
        throw error;
    }
}
