import { useState, useEffect } from 'react'
import { getAllDoctors, addDoctor, editDoctor, deleteDoctor } from '../../api/admin'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'

const empty = { name: '', email: '', p_no: '', password: '', specialization: '' }

export default function AdminDoctors() {
  const [doctors,  setDoctors]  = useState([])
  const [loading,  setLoading]  = useState(true)
  const [modal,    setModal]    = useState(null) // 'add' | 'edit' | 'delete'
  const [selected, setSelected] = useState(null)
  const [form,     setForm]     = useState(empty)
  const [alert,    setAlert]    = useState({ type: '', message: '' })
  const [saving,   setSaving]   = useState(false)

  const fetchDoctors = async () => {
    try {
      const res = await getAllDoctors()
      setDoctors(res.data.doctors || [])
    } catch { setDoctors([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchDoctors() }, [])

  const showAlert = (type, message) => {
    setAlert({ type, message })
    setTimeout(() => setAlert({ type: '', message: '' }), 4000)
  }

  const openAdd  = () => { setForm(empty); setModal('add') }
  const openEdit = (doc) => { setForm({ name: doc.name, email: doc.email, p_no: doc.p_no, password: '', specialization: doc.specialization || '' }); setSelected(doc); setModal('edit') }
  const openDelete = (doc) => { setSelected(doc); setModal('delete') }
  const closeModal = () => { setModal(null); setSelected(null) }

  const handleSave = async () => {
    if (!form.name || !form.email || !form.p_no) return showAlert('error', 'Name, email and phone are required')
    if (modal === 'add' && !form.password)       return showAlert('error', 'Password is required')
    setSaving(true)
    try {
      if (modal === 'add') {
        await addDoctor(form)
        showAlert('success', 'Doctor added successfully')
      } else {
        const payload = { name: form.name, p_no: form.p_no, specialization: form.specialization }
        await editDoctor(selected.email, payload)
        showAlert('success', 'Doctor updated successfully')
      }
      closeModal()
      fetchDoctors()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Operation failed')
    } finally { setSaving(false) }
  }

  const handleDelete = async () => {
    setSaving(true)
    try {
      await deleteDoctor(selected.email)
      showAlert('success', 'Doctor deleted successfully')
      closeModal()
      fetchDoctors()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Delete failed')
    } finally { setSaving(false) }
  }

  const columns = [
    { key: 'name',           label: 'Name' },
    { key: 'email',          label: 'Email' },
    { key: 'p_no',           label: 'Phone' },
    { key: 'specialization', label: 'Specialization' },
    { key: 'actions', label: 'Actions', render: (row) => (
      <div className="flex gap-2">
        <button onClick={() => openEdit(row)}   className="text-blue-600 hover:underline text-sm">Edit</button>
        <button onClick={() => openDelete(row)} className="text-red-600 hover:underline text-sm">Delete</button>
      </div>
    )},
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Doctors"
        subtitle={`${doctors.length} registered doctors`}
        action={<button onClick={openAdd} className="btn-primary">+ Add Doctor</button>}
      />

      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />
      <Table columns={columns} data={doctors} emptyMessage="No doctors found" />

      {(modal === 'add' || modal === 'edit') && (
        <Modal title={modal === 'add' ? 'Add Doctor' : 'Edit Doctor'} onClose={closeModal}>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
              <input className="input-field" value={form.name} onChange={e => setForm({...form, name: e.target.value})} placeholder="Dr. Ali Hassan" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input className="input-field" value={form.email} onChange={e => setForm({...form, email: e.target.value})} placeholder="ali@hospital.com" disabled={modal === 'edit'} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
              <input className="input-field" value={form.p_no} onChange={e => setForm({...form, p_no: e.target.value})} placeholder="03001234567" />
            </div>
            {modal === 'add' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Password</label>
                <input type="password" className="input-field" value={form.password} onChange={e => setForm({...form, password: e.target.value})} placeholder="••••••••" />
              </div>
            )}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Specialization</label>
              <input className="input-field" value={form.specialization} onChange={e => setForm({...form, specialization: e.target.value})} placeholder="Cardiology" />
            </div>
            <div className="flex gap-3 pt-2">
              <button onClick={handleSave} disabled={saving} className="btn-primary flex-1">{saving ? 'Saving...' : 'Save'}</button>
              <button onClick={closeModal} className="btn-secondary flex-1">Cancel</button>
            </div>
          </div>
        </Modal>
      )}

      {modal === 'delete' && (
        <Modal title="Delete Doctor" onClose={closeModal}>
          <p className="text-gray-600 mb-6">Are you sure you want to delete <strong>{selected?.name}</strong>? This cannot be undone.</p>
          <div className="flex gap-3">
            <button onClick={handleDelete} disabled={saving} className="btn-danger flex-1">{saving ? 'Deleting...' : 'Delete'}</button>
            <button onClick={closeModal} className="btn-secondary flex-1">Cancel</button>
          </div>
        </Modal>
      )}
    </div>
  )
}
