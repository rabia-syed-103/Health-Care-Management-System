import { useState, useEffect } from 'react'
import { getAllDonors, addDonor, editDonor, deleteDonor } from '../../api/donors'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'

const empty = { name: '', email: '', blood_group: '', p_no: '' }
const bloodGroups = ['A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-']

export default function BloodManagerDonors() {
  const [donors,   setDonors]   = useState([])
  const [filtered, setFiltered] = useState([])
  const [loading,  setLoading]  = useState(true)
  const [modal,    setModal]    = useState(null)
  const [selected, setSelected] = useState(null)
  const [form,     setForm]     = useState(empty)
  const [alert,    setAlert]    = useState({ type: '', message: '' })
  const [saving,   setSaving]   = useState(false)
  const [search,   setSearch]   = useState('')
  const [filterBG, setFilterBG] = useState('')

  const fetchData = async () => {
    try {
      const res = await getAllDonors()
      const data = res.data.donors || []
      setDonors(data)
      setFiltered(data)
    } catch { setDonors([]); setFiltered([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  useEffect(() => {
    let result = donors
    if (search)   result = result.filter(d => d.name?.toLowerCase().includes(search.toLowerCase()) || d.email?.toLowerCase().includes(search.toLowerCase()))
    if (filterBG) result = result.filter(d => d.blood_group?.trim() === filterBG)
    setFiltered(result)
  }, [search, filterBG, donors])

  const showAlert  = (type, message) => { setAlert({ type, message }); setTimeout(() => setAlert({ type: '', message: '' }), 4000) }
  const openAdd    = () => { setForm(empty); setModal('add') }
  const openEdit   = (row) => { setForm({ name: row.name, email: row.email, blood_group: row.blood_group?.trim(), p_no: row.p_no }); setSelected(row); setModal('edit') }
  const openDelete = (row) => { setSelected(row); setModal('delete') }
  const closeModal = () => { setModal(null); setSelected(null) }

  const handleSave = async () => {
    if (!form.name || !form.email || !form.blood_group || !form.p_no) return showAlert('error', 'All fields are required')
    setSaving(true)
    try {
      if (modal === 'add') {
        await addDonor(form)
        showAlert('success', 'Donor added successfully')
      } else {
        await editDonor(selected.email, { name: form.name, p_no: form.p_no })
        showAlert('success', 'Donor updated successfully')
      }
      closeModal(); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Operation failed')
    } finally { setSaving(false) }
  }

  const handleDelete = async () => {
    setSaving(true)
    try {
      await deleteDonor(selected.email)
      showAlert('success', 'Donor deleted')
      closeModal(); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Delete failed')
    } finally { setSaving(false) }
  }

  const columns = [
    { key: 'name',         label: 'Name' },
    { key: 'email',        label: 'Email' },
    { key: 'blood_group',  label: 'Blood Group', render: (row) => <span className="badge-danger">{row.blood_group?.trim()}</span> },
    { key: 'p_no',         label: 'Phone' },
    { key: 'last_donate',  label: 'Last Donation' },
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
        title="Donors"
        subtitle={`${filtered.length} of ${donors.length} donors`}
        action={<button onClick={openAdd} className="btn-primary">+ Add Donor</button>}
      />

      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

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
          <button onClick={() => { setSearch(''); setFilterBG('') }} className="btn-secondary">Clear</button>
        )}
      </div>

      <Table columns={columns} data={filtered} emptyMessage="No donors found" />

      {(modal === 'add' || modal === 'edit') && (
        <Modal title={modal === 'add' ? 'Add Donor' : 'Edit Donor'} onClose={closeModal}>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
              <input className="input-field" value={form.name} onChange={e => setForm({...form, name: e.target.value})} placeholder="Ahmed Khan" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input className="input-field" value={form.email} onChange={e => setForm({...form, email: e.target.value})} disabled={modal === 'edit'} placeholder="ahmed@donor.com" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Blood Group</label>
              <select className="input-field" value={form.blood_group} onChange={e => setForm({...form, blood_group: e.target.value})} disabled={modal === 'edit'}>
                <option value="">Select blood group</option>
                {bloodGroups.map(bg => <option key={bg} value={bg}>{bg}</option>)}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
              <input className="input-field" value={form.p_no} onChange={e => setForm({...form, p_no: e.target.value})} placeholder="03111234567" />
            </div>
            <div className="flex gap-3 pt-2">
              <button onClick={handleSave} disabled={saving} className="btn-primary flex-1">{saving ? 'Saving...' : 'Save'}</button>
              <button onClick={closeModal} className="btn-secondary flex-1">Cancel</button>
            </div>
          </div>
        </Modal>
      )}

      {modal === 'delete' && (
        <Modal title="Delete Donor" onClose={closeModal}>
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
