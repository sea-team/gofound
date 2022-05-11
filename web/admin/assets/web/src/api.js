import axios from 'axios'

const BASE_URL = process.env.NODE_ENV === 'production' ? '/api' : 'http://127.0.0.1:5678/api'

function request(url, method = 'get', data) {
  return axios({
    baseURL: BASE_URL,
    url: url,
    method: method,
    data: data,
  })
}

export default {
  getDatabase() {
    return request('/db/list')
  },
  query(db, params) {

    return request(`/query?database=${db}`, 'post', {
      ...params,
      highlight: params.highlight ? {
        preTag: '<em style=\'color:red\'>',
        postTag: '</em>',
      } : null,
    })
  },
  remove(db, id) {
    return request(`/remove?database=${db}`, 'post', { id })
  },
  gc() {
    return request('/gc')
  },
  getStatus() {
    return request('/status')
  },
  addIndex(db, index) {
    return request(`/index?database=${db}`, 'post', index )
  },
  drop(db){
    return request(`/db/drop?database=${db}`)
  },
  create(db){
    return request(`/db/create?database=${db}`)
  }
}
