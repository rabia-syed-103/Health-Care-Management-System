import api from './axios'

export const getAllDonors = ()            => api.get('/donors/')
export const getDonor    = (email)       => api.get(`/donors/${email}`)
export const addDonor    = (data)        => api.post('/donors/', data)
export const editDonor   = (email, data) => api.put(`/donors/${email}`, data)
export const deleteDonor = (email)       => api.delete(`/donors/${email}`)