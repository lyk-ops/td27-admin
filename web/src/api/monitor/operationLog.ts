import { request } from "@/utils/service"

export interface OrData {
  ID: number
  createdAt: string
  ip: string
  method: string
  path: string
  status: number
  userAgent: string
  reqParam: string
  respData: string
  respTime: number
  userName: string
}

export interface OrDataPageInfo {
  list: OrData[]
  total: number
  page: number
  pageSize: number
}

interface reqOrList extends PageInfo {
  path?: string
  method?: string
  status?: number
  asc?: boolean
}

// 分页获取操作记录
export function getOrListApi(data: reqOrList) {
  return request<ApiResponseData<OrDataPageInfo>>({
    url: "/opl/getOplList",
    method: "post",
    data
  })
}

// 删除操作记录
export function deleteOrApi(data: CId) {
  return request<ApiResponseData<null>>({
    url: "/opl/deleteOpl",
    method: "post",
    data
  })
}

// 批量删除操作记录
export function deleteOrByIdsApi(data: CIds) {
  return request<ApiResponseData<null>>({
    url: "/opl/deleteOplByIds",
    method: "post",
    data
  })
}
