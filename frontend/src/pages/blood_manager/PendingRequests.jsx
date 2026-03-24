import { useState, useEffect } from 'react'
import { getPendingRequests, fulfillBloodRequest } from '../../api/blood'
import { useAuth } from '../../context/AuthContext'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'

export default function BloodManagerRequests() {
  const { user } = useAuth()
  const [requests,  setRequests]  = useState([])
  const [loading,   setLoading]   = useState(true)
  const [selected,  setSelected]  = useState(null)
  const [modal,     setModal]     = useState(false)
  const [alert,     setAlert]     = useState({ type: '', message: '' })
  const [saving,    setSaving]    = useState(false)
  const [result,    setResult]    = useState(null)
  const [managerId, setManagerId] = useState(user?.user_id || '')

  const fetchData = async () => {
    try {
      const res = await getPendingRequests()
      setRequests(res.data.pending_requests || [])
    } catch { setRequests([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  const showAlert = (type, message) => { setAlert({ type, message }); setTimeout(() => setAlert({ type: '', message: '' }), 5000) }

  const openFulfill = (row) => { setSelected(row); setResult(null); setModal(true) }

  const handleFulfill = async () => {
    if (!managerId) return showAlert('error', 'Manager ID is required')
    setSaving(true)
    try {
      const res = await fulfillBloodRequest({
        request_id: selected.request_id,
        manager_id: parseInt(managerId),
      })
      setResult(res.data)
      showAlert('success', 'Blood request fulfilled successfully')
      fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Fulfillment failed — transaction rolled back')
    } finally { setSaving(false) }
  }

  const columns = [
    { key: 'request_id',    label: 'Request ID' },
    { key: 'patient_name',  label: 'Patient' },
    { key: 'blood_group',   label: 'Blood Group', render: (row) => <span className="badge-danger">{row.blood_group?.trim()}</span> },
    { key: 'quantity_needed', label: 'Units Needed' },
    { key: 'doctor_name',   label: 'Requested By' },
    { key: 'created_at',    label: 'Date' },
    { key: 'actions', label: 'Actions', render: (row) => (
      <button onClick={() => openFulfill(row)} className="btn-primary text-xs py-1 px-3">
        Fulfill
      </button>
    )},
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Pending Blood Requests"
        subtitle={`${requests.length} requests awaiting fulfillment`}
      />
      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />
      <Table columns={columns} data={requests} emptyMessage="No pending requests" />

      {modal && (
        <Modal title="Fulfill Blood Request" onClose={() => { setModal(false); setResult(null) }}>
          {!result ? (
            <div className="space-y-4">
              <div className="p-3 bg-gray-50 rounded-lg text-sm space-y-1">
                <p><span className="text-gray-500">Patient:</span> <strong>{selected?.patient_name}</strong></p>
                <p><span className="text-gray-500">Blood Group:</span> <strong>{selected?.blood_group?.trim()}</strong></p>
                <p><span className="text-gray-500">Units Needed:</span> <strong>{selected?.quantity_needed}</strong></p>
                <p><span className="text-gray-500">Requested By:</span> <strong>{selected?.doctor_name}</strong></p>
              </div>
              <div className="p-3 bg-blue-50 rounded-lg text-sm text-blue-700 border border-blue-200">
                ℹ️ The system will automatically select the best compatible blood — oldest units first.
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Manager ID</label>
                <input
                  type="number"
                  className="input-field"
                  value={managerId}
                  onChange={e => setManagerId(e.target.value)}
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button onClick={handleFulfill} disabled={saving} className="btn-primary flex-1">
                  {saving ? 'Processing...' : 'Confirm Fulfillment'}
                </button>
                <button onClick={() => setModal(false)} className="btn-secondary flex-1">Cancel</button>
              </div>
            </div>
          ) : (
            <div className="space-y-3">
              <div className="p-4 bg-green-50 rounded-lg border border-green-200">
                <h3 className="font-semibold text-green-800 mb-3">Request Fulfilled Successfully</h3>
                <div className="space-y-1 text-sm text-green-700">
                  <p>Request ID: <strong>#{result.request_id}</strong></p>
                  <p>Patient: <strong>{result.patient}</strong></p>
                  <p>Blood Group Used: <strong>{result.patient_blood_gr}</strong></p>
                  <p>Blood ID Used: <strong>#{result.blood_id_used}</strong></p>
                  <p>Units Provided: <strong>{result.units_provided}</strong></p>
                  <p>Blood Expiry: <strong>{result.blood_expiry}</strong></p>
                  <p>Manager: <strong>{result.manager}</strong></p>
                </div>
              </div>
              <button onClick={() => { setModal(false); setResult(null) }} className="btn-primary w-full">Done</button>
            </div>
          )}
        </Modal>
      )}
    </div>
  )
}
