import { useState, useEffect } from 'react'
import { getAllERPatients, addERPatient } from '../../api/er'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'

const empty = { name: '', age: '', p_no: '', arrival_time: '' }

export default function ReceptionistERPatients() {
  const [data,    setData]    = useState([])
  const [loading, setLoading] = useState(true)
  const [modal,   setModal]   = useState(false)
  const [form,    setForm]    = useState(empty)
  const [alert,   setAlert]   = useState({ type: '', message: '' })
  const [saving,  setSaving]  = useState(false)

  const fetchData = async () => {
    try {
      const res = await getAllERPatients()
      setData(res.data.er_patients || [])
    } catch { setData([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  const showAlert = (type, message) => { setAlert({ type, message }); setTimeout(() => setAlert({ type: '', message: '' }), 4000) }

  const handleSave = async () => {
    if (!form.name || !form.age || !form.p_no || !form.arrival_time) return showAlert('error', 'All fields are required')
    setSaving(true)
    try {
      await addERPatient({ ...form, age: parseInt(form.age) })
      showAlert('success', 'ER patient registered successfully')
      setModal(false); setForm(empty); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Failed to add ER patient')
    } finally { setSaving(false) }
  }

  const columns = [
    { key: 'id',           label: 'ID' },
    { key: 'name',         label: 'Name' },
    { key: 'age',          label: 'Age' },
    { key: 'p_no',         label: 'Phone' },
    { key: 'arrival_time', label: 'Arrival Time' },
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="ER Patients"
        subtitle={`${data.length} ER walk-in patients`}
        action={<button onClick={() => { setForm(empty); setModal(true) }} className="btn-primary">+ Register ER Patient</button>}
      />
      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />
      <Table columns={columns} data={data} emptyMessage="No ER patients found" />

      {modal && (
        <Modal title="Register ER Patient" onClose={() => setModal(false)}>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
              <input className="input-field" value={form.name} onChange={e => setForm({...form, name: e.target.value})} placeholder="Ahmed Khan" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Age</label>
              <input type="number" className="input-field" value={form.age} onChange={e => setForm({...form, age: e.target.value})} placeholder="35" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
              <input className="input-field" value={form.p_no} onChange={e => setForm({...form, p_no: e.target.value})} placeholder="03001234567" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Arrival Time</label>
              <input type="time" className="input-field" value={form.arrival_time} onChange={e => setForm({...form, arrival_time: e.target.value})} />
            </div>
            <div className="flex gap-3 pt-2">
              <button onClick={handleSave} disabled={saving} className="btn-primary flex-1">{saving ? 'Saving...' : 'Register'}</button>
              <button onClick={() => setModal(false)} className="btn-secondary flex-1">Cancel</button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  )
}
