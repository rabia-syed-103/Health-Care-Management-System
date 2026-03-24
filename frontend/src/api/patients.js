import api from './axios'

export const getAllPatients  = ()           => api.get('/patients/')
export const getPatient      = (email)      => api.get(`/patients/${email}`)
export const addPatient      = (data)       => api.post('/patients/', data)
export const editPatient     = (email, data)=> api.put(`/patients/${email}`, data)
export const deletePatient   = (email)      => api.delete(`/patients/${email}`)