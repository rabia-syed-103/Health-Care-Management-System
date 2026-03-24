import api from './axios'

export const getAllDoctors    = ()            => api.get('/admin/doctors')
export const addDoctor        = (data)        => api.post('/admin/doctors', data)
export const editDoctor       = (email, data) => api.put(`/admin/doctors/${email}`, data)
export const deleteDoctor     = (email)       => api.delete(`/admin/doctors/${email}`)

export const getAllReceptionists = ()            => api.get('/admin/receptionists')
export const addReceptionist     = (data)        => api.post('/admin/receptionists', data)
export const editReceptionist    = (email, data) => api.put(`/admin/receptionists/${email}`, data)
export const deleteReceptionist  = (email)       => api.delete(`/admin/receptionists/${email}`)

export const getAllPharmacists = ()            => api.get('/admin/pharmacists')
export const addPharmacist     = (data)        => api.post('/admin/pharmacists', data)
export const editPharmacist    = (email, data) => api.put(`/admin/pharmacists/${email}`, data)
export const deletePharmacist  = (email)       => api.delete(`/admin/pharmacists/${email}`)

export const getAllBloodManagers = ()            => api.get('/admin/blood-managers')
export const addBloodManager     = (data)        => api.post('/admin/blood-managers', data)
export const editBloodManager    = (email, data) => api.put(`/admin/blood-managers/${email}`, data)
export const deleteBloodManager  = (email)       => api.delete(`/admin/blood-managers/${email}`)