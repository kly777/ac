import { ref, watch } from "vue";

// WebSocket服务类
class WebSocketService {
    private socket: WebSocket | null = null;
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectInterval = 3000; // 3秒重连间隔
    private url: string;

    // 状态
    isConnected = ref(false);
    messages = ref<any[]>([]);
    connectionStatus = ref('disconnected'); // 连接状态: 'connecting', 'connected', 'disconnected'

    constructor(url: string) {
        this.url = url;
    }

    // 连接WebSocket
    connect(): void {
        this.connectionStatus.value = 'connecting';
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            this.isConnected.value = true;
            this.reconnectAttempts = 0;
            this.connectionStatus.value = 'connected';
            console.log("WebSocket连接已建立");
        };

        this.socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.messages.value.push(data);
            } catch (error) {
                console.error("解析WebSocket消息失败:", error);
            }
        };

        this.socket.onerror = (error) => {
            console.error("WebSocket错误:", error);
        };

        this.socket.onclose = () => {
            this.isConnected.value = false;
            this.connectionStatus.value = 'disconnected';
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
    send(message: any): void {
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

    // 监听特定类型消息
    onMessageType(type: string, callback: (data: any) => void): void {
        const unwatch = watch(
            this.messages,
            () => {
                const newMessages = this.messages.value.filter(
                    (msg) => msg.type === type
                );
                newMessages.forEach((msg) => callback(msg.data));
                // 移除已处理的消息
                this.messages.value = this.messages.value.filter(
                    (msg) => !newMessages.includes(msg)
                );
            },
            { deep: true }
        );
    }
    // 获取当前连接状态
    getStatus() {
        return this.connectionStatus.value;
    }
}

// 创建单例实例
const wsService = new WebSocketService("ws://localhost:8080/ws");

export default wsService;
