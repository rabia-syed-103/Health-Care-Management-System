import { useState, useEffect } from 'react'
import { getAllPatients, addPatient, editPatient, deletePatient } from '../../api/patients'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'
import { exportPDF } from '../../utils/pdfExport'

const empty = { name: '', email: '', b_gr: '', p_no: '' }
const bloodGroups = ['A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-']

export default function ReceptionistPatients() {
  const [patients,  setPatients]  = useState([])
  const [filtered,  setFiltered]  = useState([])
  const [loading,   setLoading]   = useState(true)
  const [modal,     setModal]     = useState(null)
  const [selected,  setSelected]  = useState(null)
  const [form,      setForm]      = useState(empty)
  const [alert,     setAlert]     = useState({ type: '', message: '' })
  const [saving,    setSaving]    = useState(false)
  const [search,    setSearch]    = useState('')
  const [filterBG,  setFilterBG]  = useState('')

  const fetchData = async () => {
    try {
      const res = await getAllPatients()
      const data = res.data.patients || []
      setPatients(data)
      setFiltered(data)
    } catch { setPatients([]); setFiltered([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  useEffect(() => {
    let result = patients
    if (search)   result = result.filter(p => p.name?.toLowerCase().includes(search.toLowerCase()) || p.email?.toLowerCase().includes(search.toLowerCase()))
    if (filterBG) result = result.filter(p => p.b_gr?.trim() === filterBG)
    setFiltered(result)
  }, [search, filterBG, patients])

  const showAlert  = (type, message) => { setAlert({ type, message }); setTimeout(() => setAlert({ type: '', message: '' }), 4000) }
  const openAdd    = () => { setForm(empty); setModal('add') }
  const openEdit   = (row) => { setForm({ name: row.name, email: row.email, b_gr: row.b_gr?.trim(), p_no: row.p_no }); setSelected(row); setModal('edit') }
  const openDelete = (row) => { setSelected(row); setModal('delete') }
  const closeModal = () => { setModal(null); setSelected(null) }

  const handleSave = async () => {
    if (!form.name || !form.email || !form.b_gr || !form.p_no) return showAlert('error', 'All fields are required')
    setSaving(true)
    try {
      if (modal === 'add') {
        await addPatient(form)
        showAlert('success', 'Patient added successfully')
      } else {
        await editPatient(selected.email, { name: form.name, b_gr: form.b_gr, p_no: form.p_no })
        showAlert('success', 'Patient updated successfully')
      }
      closeModal(); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Operation failed')
    } finally { setSaving(false) }
  }

  const handleDelete = async () => {
    setSaving(true)
    try {
      await deletePatient(selected.email)
      showAlert('success', 'Patient deleted')
      closeModal(); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Delete failed')
    } finally { setSaving(false) }
  }
   const downloadPDF = () => exportPDF(
    'Patients Report',
    ['Name', 'Email', 'Blood Group', 'Phone'],
    filtered.map(r => [r.name, r.email, r.b_gr?.trim(), r.p_no]),
    'patients-report'
  )
  const columns = [
    { key: 'name',  label: 'Name' },
    { key: 'email', label: 'Email' },
    { key: 'b_gr',  label: 'Blood Group', render: (row) => <span className="badge-info">{row.b_gr?.trim()}</span> },
    { key: 'p_no',  label: 'Phone' },
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
      title="Patients"
      subtitle={`${filtered.length} of ${patients.length} patients`}
      action={
        <div className="flex gap-2">
          <button onClick={downloadPDF} className="btn-secondary">⬇ Download PDF</button>
          <button onClick={openAdd} className="btn-primary">+ Add Patient</button>
        </div>
      }
    />

      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

      {/* Search & Filter */}
      <div className="flex gap-3 mb-4">
        <input
          className="input-field max-w-sm"
          placeholder="Search by name or email..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <select
          className="input-field max-w-xs"
          value={filterBG}
          onChange={e => setFilterBG(e.target.value)}
        >
          <option value="">All Blood Groups</option>
          {bloodGroups.map(bg => <option key={bg} value={bg}>{bg}</option>)}
        </select>
        {(search || filterBG) && (
          <button onClick={() => { setSearch(''); setFilterBG('') }} className="btn-secondary">
            Clear
          </button>
        )}
      </div>

      <Table columns={columns} data={filtered} emptyMessage="No patients found" />

      {(modal === 'add' || modal === 'edit') && (
        <Modal title={modal === 'add' ? 'Add Patient' : 'Edit Patient'} onClose={closeModal}>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
              <input className="input-field" value={form.name} onChange={e => setForm({...form, name: e.target.value})} placeholder="John Doe" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input className="input-field" value={form.email} onChange={e => setForm({...form, email: e.target.value})} disabled={modal === 'edit'} placeholder="john@email.com" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Blood Group</label>
              <select className="input-field" value={form.b_gr} onChange={e => setForm({...form, b_gr: e.target.value})}>
                <option value="">Select blood group</option>
                {bloodGroups.map(bg => <option key={bg} value={bg}>{bg}</option>)}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
              <input className="input-field" value={form.p_no} onChange={e => setForm({...form, p_no: e.target.value})} placeholder="03001234567" />
            </div>
            <div className="flex gap-3 pt-2">
              <button onClick={handleSave} disabled={saving} className="btn-primary flex-1">{saving ? 'Saving...' : 'Save'}</button>
              <button onClick={closeModal} className="btn-secondary flex-1">Cancel</button>
            </div>
          </div>
        </Modal>
      )}

      {modal === 'delete' && (
        <Modal title="Delete Patient" onClose={closeModal}>
          <p className="text-gray-600 mb-6">Are you sure you want to delete <strong>{selected?.name}</strong>?</p>
          <div className="flex gap-3">
            <button onClick={handleDelete} disabled={saving} className="btn-danger flex-1">{saving ? 'Deleting...' : 'Delete'}</button>
            <button onClick={closeModal} className="btn-secondary flex-1">Cancel</button>
          </div>
        </Modal>
      )}
    </div>
  )
}
