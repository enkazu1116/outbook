/**
 * Backend API クライアント
 * 
 * 環境変数 NEXT_PUBLIC_API_BASE_URL で backend の URL を設定できます。
 * デフォルトは http://localhost:8080 です。
 */

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public response?: Response
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * APIリクエストのオプション
 */
interface RequestOptions extends RequestInit {
  /** 認証トークン（Bearerトークン） */
  token?: string;
  /** リクエストボディ（オブジェクトは自動的にJSONに変換） */
  body?: unknown;
}

/**
 * 汎用APIリクエスト関数
 */
async function request<T>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<T> {
  const { token, body, headers, ...fetchOptions } = options;

  // URLを構築
  const url = endpoint.startsWith('http') 
    ? endpoint 
    : `${API_BASE_URL}${endpoint}`;

  // ヘッダーを構築
  const requestHeaders: HeadersInit = {
    'Content-Type': 'application/json',
    ...headers,
  };

  // 認証トークンが提供されている場合は追加
  if (token) {
    requestHeaders['Authorization'] = `Bearer ${token}`;
  }

  // ボディをJSONに変換（オブジェクトの場合）
  let requestBody: string | undefined;
  if (body !== undefined) {
    if (typeof body === 'string') {
      requestBody = body;
    } else {
      requestBody = JSON.stringify(body);
    }
  }

  try {
    const response = await fetch(url, {
      ...fetchOptions,
      headers: requestHeaders,
      body: requestBody,
    });

    // エラーレスポンスを処理
    if (!response.ok) {
      let errorMessage = `HTTP Error: ${response.status} ${response.statusText}`;
      try {
        const errorData = await response.json();
        errorMessage = errorData.message || errorData.error || errorMessage;
      } catch {
        // JSONパースに失敗した場合はデフォルトメッセージを使用
      }
      throw new ApiError(errorMessage, response.status, response);
    }

    // レスポンスが空の場合は null を返す
    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      return null as T;
    }

    return await response.json();
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError(
      error instanceof Error ? error.message : 'ネットワークエラーが発生しました',
      0
    );
  }
}

/**
 * GETリクエスト
 */
export function get<T>(endpoint: string, options?: RequestOptions): Promise<T> {
  return request<T>(endpoint, { ...options, method: 'GET' });
}

/**
 * POSTリクエスト
 */
export function post<T>(
  endpoint: string,
  body?: unknown,
  options?: RequestOptions
): Promise<T> {
  return request<T>(endpoint, { ...options, method: 'POST', body });
}

/**
 * PUTリクエスト
 */
export function put<T>(
  endpoint: string,
  body?: unknown,
  options?: RequestOptions
): Promise<T> {
  return request<T>(endpoint, { ...options, method: 'PUT', body });
}

/**
 * DELETEリクエスト
 */
export function del<T>(endpoint: string, options?: RequestOptions): Promise<T> {
  return request<T>(endpoint, { ...options, method: 'DELETE' });
}

/**
 * PATCHリクエスト
 */
export function patch<T>(
  endpoint: string,
  body?: unknown,
  options?: RequestOptions
): Promise<T> {
  return request<T>(endpoint, { ...options, method: 'PATCH', body });
}

/**
 * ヘルスチェック
 */
export async function healthCheck(): Promise<boolean> {
  try {
    const response = await get<{ status?: string }>('/healthz');
    return response?.status === 'OK' || true;
  } catch {
    return false;
  }
}

/**
 * ユーザー関連API
 */
export const userApi = {
  /**
   * ユーザー登録
   */
  register: (email: string, password: string, name: string) =>
    post<{ message: string }>('/api/users', { email, password, name }),

  /**
   * ユーザー情報取得
   */
  getUser: (id: string) =>
    get<{ id: number }>(`/api/users?id=${id}`),
};

