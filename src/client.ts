import type {
  YapiConfig,
  YapiApiResponse,
  YapiProject,
  YapiCategory,
  YapiInterfaceDetail,
  CreateInterfaceParams,
  UpdateInterfaceParams,
  YapiLoginResponse,
  CreateCategoryParams,
  YapiCategoryCreated,
} from "./types.js";

export class YapiClient {
  private baseUrl: string;
  private cookies: string | null = null;
  private loginCredentials: { email: string; password: string };

  constructor(config: YapiConfig) {
    this.baseUrl = config.baseUrl.replace(/\/$/, "");
    this.loginCredentials = { email: config.email, password: config.password };
  }

  /**
   * 通过邮箱和密码登录，获取 cookie
   */
  async login(email: string, password: string): Promise<void> {
    const url = new URL(`${this.baseUrl}/api/user/login`);
    const response = await fetch(url.toString(), {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    // 使用 getSetCookie() 方法获取所有 set-cookie 头
    const setCookies = response.headers.getSetCookie?.() || [];
    
    if (setCookies.length > 0) {
      // 解析 set-cookie 头，提取 _yapi_token 和 _yapi_uid
      const cookies: string[] = [];
      
      for (const cookieStr of setCookies) {
        // 提取 cookie 名称和值
        const match = cookieStr.match(/(_yapi_\w+)=([^;]+)/);
        if (match) {
          cookies.push(`${match[1]}=${match[2]}`);
        }
      }
      
      if (cookies.length > 0) {
        this.cookies = cookies.join("; ");
        console.error(`[yapi-mcp] Login successful, cookies: ${cookies.length} found`);
        return;
      }
    }

    // 备用方案：尝试从 get("set-cookie") 获取
    const setCookieHeader = response.headers.get("set-cookie");
    if (setCookieHeader) {
      const cookies: string[] = [];
      const cookiePairs = setCookieHeader.split(", ");
      
      for (const pair of cookiePairs) {
        const match = pair.match(/(_yapi_\w+)=([^;]+)/);
        if (match) {
          cookies.push(`${match[1]}=${match[2]}`);
        }
      }
      
      if (cookies.length > 0) {
        this.cookies = cookies.join("; ");
        console.error(`[yapi-mcp] Login successful (fallback), cookies: ${cookies.length} found`);
        return;
      }
    }

    // 检查响应是否成功
    const data = (await response.clone().json()) as YapiApiResponse<YapiLoginResponse>;
    if (data.errcode !== 0) {
      throw new Error(`YAPI Login Error: ${data.errmsg} (code: ${data.errcode})`);
    }

    throw new Error("YAPI Login failed: No cookies received");
  }

  /**
   * 确保已登录（有 cookies）
   */
  private async ensureAuthenticated(): Promise<void> {
    if (this.cookies) {
      return;
    }

    await this.login(this.loginCredentials.email, this.loginCredentials.password);
  }

  private async request<T>(
    path: string,
    options: {
      method?: string;
      body?: Record<string, unknown>;
    } = {}
  ): Promise<YapiApiResponse<T>> {
    await this.ensureAuthenticated();

    const url = new URL(`${this.baseUrl}${path}`);

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      "Cookie": this.cookies!,
    };

    const fetchOptions: RequestInit = {
      method: options.method || "GET",
      headers,
    };

    if (options.body) {
      fetchOptions.body = JSON.stringify(options.body);
    }

    const response = await fetch(url.toString(), fetchOptions);
    const data = (await response.json()) as YapiApiResponse<T>;

    if (data.errcode !== 0) {
      throw new Error(`YAPI Error: ${data.errmsg} (code: ${data.errcode})`);
    }

    return data;
  }

  /**
   * 获取项目信息
   * @param projectId 项目ID
   */
  async getProject(projectId: number): Promise<YapiProject> {
    const response = await this.request<YapiProject>(
      `/api/project/get?id=${projectId}`
    );
    return response.data;
  }

  /**
   * 获取项目列表（通过 token 获取可访问的项目）
   */
  async getProjectList(): Promise<YapiProject[]> {
    const response = await this.request<YapiProject[]>("/api/project/list");
    return response.data;
  }

  /**
   * 获取项目的分类列表
   * @param projectId 项目ID
   */
  async getCategories(projectId: number): Promise<YapiCategory[]> {
    const response = await this.request<YapiCategory[]>(
      `/api/interface/getCatMenu?project_id=${projectId}`
    );
    return response.data;
  }

  /**
   * 获取分类下的接口列表
   * @param catid 分类ID
   */
  async getInterfacesByCategory(catid: number): Promise<YapiCategory> {
    const response = await this.request<YapiCategory>(
      `/api/interface/list_cat?catid=${catid}`
    );
    return response.data;
  }

  /**
   * 获取项目的所有接口列表
   * @param projectId 项目ID
   */
  async getInterfacesByProject(projectId: number): Promise<YapiInterfaceDetail[]> {
    const response = await this.request<YapiInterfaceDetail[]>(
      `/api/interface/list?project_id=${projectId}`
    );
    return response.data;
  }

  /**
   * 获取接口详情
   * @param id 接口ID
   */
  async getInterface(id: number): Promise<YapiInterfaceDetail> {
    const response = await this.request<YapiInterfaceDetail>(
      `/api/interface/get?id=${id}`
    );
    return response.data;
  }

  /**
   * 新增接口
   * @param params 接口参数
   */
  async createInterface(params: CreateInterfaceParams): Promise<YapiInterfaceDetail> {
    const response = await this.request<YapiInterfaceDetail>(
      "/api/interface/save",
      {
        method: "POST",
        body: params as unknown as Record<string, unknown>,
      }
    );
    return response.data;
  }

  /**
   * 修改接口
   * @param params 接口参数
   */
  async updateInterface(params: UpdateInterfaceParams): Promise<YapiInterfaceDetail> {
    const response = await this.request<YapiInterfaceDetail>(
      "/api/interface/up",
      {
        method: "POST",
        body: params as unknown as Record<string, unknown>,
      }
    );
    return response.data;
  }

  /**
   * 新增分类
   * @param params 分类参数
   */
  async createCategory(params: CreateCategoryParams): Promise<YapiCategoryCreated> {
    const response = await this.request<YapiCategoryCreated>(
      "/api/interface/add_cat",
      {
        method: "POST",
        body: params as unknown as Record<string, unknown>,
      }
    );
    return response.data;
  }
}
