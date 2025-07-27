/**
 * 前后端交互类型定义
 */

// 统一WebSocket消息结构
export interface WSMessage<T = any> {
  type: string;
  data: T;
}

// 问题接口请求类型
export interface QuestionRequest {
  content: string;
}



// 连接状态消息数据结构
export interface ConnectionStatusData {
  status: 'connected' | 'disconnected' | 'connecting';
}

// 实时数据消息数据结构
export interface RealTimeData {
  // 根据实际数据结构定义
  [key: string]: any;
}

// 应用消息类型
export type AppMessageType =
    | "ai_response"
    | "informer"
    | "error";