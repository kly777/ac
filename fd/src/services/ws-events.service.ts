import wsService from './websocket.service';




// 注册AI响应回调

type AIResponseData = {
  status: 'ing' | 'done';
  content: string;
}
export function onAIResponse(callback: (content: AIResponseData) => void): void {
    wsService.on("ai_response", (data: AIResponseData) => {
        callback(data);
    });
}




