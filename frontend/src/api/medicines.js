import api from './axios'

export const getMedicineStock  = ()      => api.get('/pharmacist/medicine-stock')
export const getMedicineByName = (name)  => api.get(`/pharmacist/medicine/${name}`)
export const addMedicine       = (data)  => api.post('/pharmacist/medicine', data)