import { ref } from "vue";
import type { AppMessageType } from "@/types";

// WebSocket事件监听器类型
type EventListener = (data: any) => void;

// 连接状态类型
type ConnectionStatus = "disconnected" | "connecting" | "connected" | "error";

// WebSocket服务类
class WebSocketService {
    private socket: WebSocket | null = null;
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectInterval = 3000; // 3秒重连间隔
    private url: string;
    private eventListeners: Map<AppMessageType, EventListener[]> = new Map();

    // 状态
    isConnected = ref(false);
    connectionStatus = ref<ConnectionStatus>("disconnected");

    constructor(url: string) {
        this.url = url;
        this.connect(); // 在构造时自动连接
    }

    // 连接WebSocket
    private connect(): void {
        this.connectionStatus.value = "connecting";
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            this.isConnected.value = true;
            this.reconnectAttempts = 0;
            this.connectionStatus.value = "connected";
            console.log("WebSocket连接已建立");
        };

        this.socket.onmessage = (event) => {
            try {
                const message: { type: AppMessageType; data: any } = JSON.parse(
                    event.data
                );
                console.log("收到WebSocket消息：", message);

                // 按消息类型分发事件
                if (message.type && this.eventListeners.has(message.type)) {
                    console.log("分发事件：", message.type);
                    const listeners =
                        this.eventListeners.get(message.type) || [];
                    listeners.forEach((listener) => listener(message.data));
                }
            } catch (error) {
                console.error("解析WebSocket消息失败:", error);
            }
        };

        this.socket.onerror = (error) => {
            console.error("WebSocket错误:", error);
            this.connectionStatus.value = "error";
        };

        this.socket.onclose = () => {
            this.isConnected.value = false;
            this.connectionStatus.value = "disconnected";
            console.log("WebSocket连接已关闭");

            // 自动重连逻辑
            if (this.reconnectAttempts < this.maxReconnectAttempts) {
                setTimeout(() => {
                    this.reconnectAttempts++;
                    console.log(
                        `尝试重新连接 (${this.reconnectAttempts}/${this.maxReconnectAttempts})`
                    );
                    this.connect();
                }, this.reconnectInterval);
            }
        };
    }

    // 发送消息
    send(message: { type: AppMessageType; data: any }): void {
        if (this.isConnected.value && this.socket) {
            this.socket.send(JSON.stringify(message));
        } else {
            console.warn("无法发送消息：WebSocket未连接");
        }
    }

    // 关闭连接
    disconnect(): void {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    // 注册事件监听器
    on(type: AppMessageType, callback: EventListener): void {
        if (!this.eventListeners.has(type)) {
            this.eventListeners.set(type, []);
        }
        this.eventListeners.get(type)?.push(callback);
    }

    // 移除事件监听器
    off(type: AppMessageType, callback: EventListener): void {
        const listeners = this.eventListeners.get(type);
        if (listeners) {
            const index = listeners.indexOf(callback);
            if (index !== -1) {
                listeners.splice(index, 1);
            }
        }
    }

    // 获取当前连接状态
    getStatus(): ConnectionStatus {
        return this.connectionStatus.value;
    }
}

// 创建单例实例
const wsService = new WebSocketService("ws://localhost:8080/ws");

export default wsService;
