/**
 * API服务封装
 */

import { request } from './http.service';

/**
 * 发送问题到后端
 * @param content 问题内容
 */
export async function Q(content: string): Promise<void> {
  await request('/q', 'POST', { content });
}