import { request } from "@/utils/service"

export interface MenusData {
  id: number
  pid: number
  name: string
  path: string
  redirect: string
  component: string
  meta: {
    hidden: boolean
    title: string
    elIcon: string
    svgIcon: string
    affix: boolean
  }
  children: MenusData[]
}

type MenusResponseData = IApiResponseData<MenusData[]>

// 获取动态路由
export function getMenus() {
  return request<MenusResponseData>({
    url: "/menu/getMenus",
    method: "get"
  })
}

export interface reqMenu {
  id: number
  pid: number
  name: string
  path: string
  redirect: string
  component: string
  meta: {
    hidden: boolean
    title: string
    icon: string
    affix: boolean
  }
}

export function addMenuApi(data: reqMenu) {
  return request<IApiResponseData<null>>({
    url: "menu/addMenu",
    method: "post",
    data
  })
}

export function editMenuApi(data: reqMenu) {
  return request<IApiResponseData<null>>({
    url: "menu/editMenu",
    method: "post",
    data
  })
}

interface reqId {
  id: number
}

export function deleteMenuApi(data: reqId) {
  return request<IApiResponseData<null>>({
    url: "menu/deleteMenu",
    method: "post",
    data
  })
}