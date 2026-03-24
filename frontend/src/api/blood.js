import api from './axios'

export const recordDonation      = (data) => api.post('/blood/donate', data)
export const createBloodRequest  = (data) => api.post('/blood-request/create', data)
export const fulfillBloodRequest = (data) => api.post('/blood-fulfill/fulfill', data)
export const getDonationHistory  = ()     => api.get('/blood-manager/donations')
export const getBloodInventory   = ()     => api.get('/blood-manager/inventory')
export const getPendingRequests  = ()     => api.get('/blood-manager/pending-requests')
export const getExpiredBlood     = ()     => api.get('/blood-manager/expired-blood')