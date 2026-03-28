import { useState, useEffect } from 'react'
import { getAllEREntries, addEREntry } from '../../api/er'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'

const empty = { er_patient_id: '', doctor_id: '', er_id: '', status: 'waiting' }
const statuses = ['waiting', 'in_treatment', 'admitted', 'discharged', 'transferred']

export default function ReceptionistEREntries() {
  const [data,    setData]    = useState([])
  const [loading, setLoading] = useState(true)
  const [modal,   setModal]   = useState(false)
  const [form,    setForm]    = useState(empty)
  const [alert,   setAlert]   = useState({ type: '', message: '' })
  const [saving,  setSaving]  = useState(false)

  const fetchData = async () => {
    try {
      const res = await getAllEREntries()
      setData(res.data.er_entries || [])
    } catch { setData([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  const showAlert = (type, message) => { setAlert({ type, message }); setTimeout(() => setAlert({ type: '', message: '' }), 4000) }

  const handleSave = async () => {
    if (!form.er_patient_id || !form.doctor_id || !form.er_id) return showAlert('error', 'All fields are required')
    setSaving(true)
    try {
      await addEREntry({
        er_patient_id: parseInt(form.er_patient_id),
        doctor_id:     parseInt(form.doctor_id),
        er_id:         parseInt(form.er_id),
        status:        form.status,
      })
      showAlert('success', 'ER entry added successfully')
      setModal(false); setForm(empty); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Failed to add ER entry')
    } finally { setSaving(false) }
  }

  const statusBadge = (status) => {
    const map = {
      waiting:      'badge-warning',
      in_treatment: 'badge-info',
      admitted:     'badge-info',
      discharged:   'badge-success',
      transferred:  'badge-danger',
    }
    return <span className={map[status] || 'badge-info'}>{status?.replace('_', ' ')}</span>
  }

  const columns = [
    { key: 'patient_name', label: 'ER Patient' },
    { key: 'doctor_name',     label: 'Doctor' },
    { key: 'id',           label: 'Shift ID' },
    {key: 'shift_date', label: 'Shift Date'},
    {key:'shift_time', label: 'Shift Time'},
    { key: 'status', label: 'Status', render: (row) => statusBadge(row.status) },
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
    <PageHeader
      title="ER Entries"
      subtitle={`${data.length} ER patient entries`}
    />
      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />
      <Table columns={columns} data={data} emptyMessage="No ER entries found" />

      {modal && (
        <Modal title="Add ER Entry" onClose={() => setModal(false)}>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">ER Patient ID</label>
              <input type="number" className="input-field" value={form.er_patient_id} onChange={e => setForm({...form, er_patient_id: e.target.value})} placeholder="1" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Doctor ID</label>
              <input type="number" className="input-field" value={form.doctor_id} onChange={e => setForm({...form, doctor_id: e.target.value})} placeholder="1" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">ER Shift ID</label>
              <input type="number" className="input-field" value={form.er_id} onChange={e => setForm({...form, er_id: e.target.value})} placeholder="1" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select className="input-field" value={form.status} onChange={e => setForm({...form, status: e.target.value})}>
                {statuses.map(s => <option key={s} value={s}>{s.replace('_', ' ')}</option>)}
              </select>
            </div>
            <div className="flex gap-3 pt-2">
              <button onClick={handleSave} disabled={saving} className="btn-primary flex-1">{saving ? 'Saving...' : 'Add Entry'}</button>
              <button onClick={() => setModal(false)} className="btn-secondary flex-1">Cancel</button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  )
}
