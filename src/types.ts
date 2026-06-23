export interface YapiConfig {
  baseUrl: string;
  email: string;
  password: string;
  loginPath?: string;
}

export interface YapiLoginResponse {
  uid: number;
  username: string;
  email: string;
  token: string;
}

export interface YapiProject {
  _id: number;
  name: string;
  basepath: string;
  desc: string;
  members: Array<{
    uid: number;
    role: string;
    username: string;
    email: string;
  }>;
}

export interface YapiCategory {
  _id: number;
  name: string;
  desc: string;
  project_id: number;
  list: YapiInterface[];
}

export interface YapiInterface {
  _id: number;
  title: string;
  path: string;
  method: string;
  project_id: number;
  catid: number;
  uid: number;
  add_time: number;
  up_time: number;
  desc: string;
  status: string;
  req_body_type: string;
  req_body_other?: string;
  req_body_form?: Array<{
    name: string;
    type: string;
    required: string;
    desc?: string;
  }>;
  req_headers?: Array<{
    name: string;
    value?: string;
    required: string;
    desc?: string;
  }>;
  req_params?: Array<{
    name: string;
    desc?: string;
  }>;
  req_query?: Array<{
    name: string;
    required: string;
    desc?: string;
  }>;
  res_body_type: string;
  res_body?: string;
}

export interface YapiInterfaceDetail extends YapiInterface {
  username: string;
  editor?: string;
}

export interface YapiApiResponse<T> {
  errcode: number;
  errmsg: string;
  data: T;
}

export interface CreateInterfaceParams {
  title: string;
  path: string;
  method: string;
  catid: number;
  project_id?: number;
  desc?: string;
  status?: string;
  req_body_type?: string;
  req_body_other?: string;
  req_body_form?: Array<{
    name: string;
    type: string;
    required: string;
    desc?: string;
  }>;
  req_headers?: Array<{
    name: string;
    value?: string;
    required: string;
    desc?: string;
  }>;
  req_params?: Array<{
    name: string;
    desc?: string;
  }>;
  req_query?: Array<{
    name: string;
    required: string;
    desc?: string;
  }>;
  res_body_type?: string;
  res_body?: string;
}

export interface UpdateInterfaceParams extends Partial<CreateInterfaceParams> {
  id: number;
}

export interface CreateCategoryParams {
  name: string;
  desc?: string;
  project_id: number;
}

export interface YapiCategoryCreated {
  _id: number;
  name: string;
  desc: string;
  project_id: number;
  uid: number;
  add_time: number;
  up_time: number;
}
