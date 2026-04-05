import { useState, useEffect } from 'react'
import { getAllPharmacists, addPharmacist, editPharmacist, deletePharmacist, getPharmacistDetail } from '../../api/admin'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'

const empty = { name: '', email: '', p_no: '', password: '' }

export default function AdminPharmacists() {
  const [data,     setData]     = useState([])
  const [loading,  setLoading]  = useState(true)
  const [modal,    setModal]    = useState(null)
  const [selected, setSelected] = useState(null)
  const [form,     setForm]     = useState(empty)
  const [alert,    setAlert]    = useState({ type: '', message: '' })
  const [saving,   setSaving]   = useState(false)
  const [detail,        setDetail]        = useState(null)
  const [detailLoading, setDetailLoading] = useState(false)
  const fetchData = async () => {
    try {
      const res = await getAllPharmacists()
      setData(res.data.pharmacists || [])
    } catch { setData([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  const showAlert  = (type, message) => { setAlert({ type, message }); setTimeout(() => setAlert({ type: '', message: '' }), 4000) }
  const openAdd    = () => { setForm(empty); setModal('add') }
  const openEdit   = (row) => { setForm({ name: row.name, email: row.email, p_no: row.p_no, password: '' }); setSelected(row); setModal('edit') }
  const openDelete = (row) => { setSelected(row); setModal('delete') }
  const closeModal = () => { setModal(null); setSelected(null) }

  const handleSave = async () => {
    if (!form.name || !form.email || !form.p_no) return showAlert('error', 'All fields are required')
    if (modal === 'add' && !form.password) return showAlert('error', 'Password is required')
    setSaving(true)
    try {
      modal === 'add'
        ? await addPharmacist(form)
        : await editPharmacist(selected.email, { name: form.name, p_no: form.p_no })
      showAlert('success', `Pharmacist ${modal === 'add' ? 'added' : 'updated'} successfully`)
      closeModal(); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Operation failed')
    } finally { setSaving(false) }
  }

  const handleDelete = async () => {
    setSaving(true)
    try {
      await deletePharmacist(selected.email)
      showAlert('success', 'Pharmacist deleted')
      closeModal(); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Delete failed')
    } finally { setSaving(false) }
  }
  const openDetail = async (row) => {
      setSelected(row)
      setModal('detail')
      setDetailLoading(true)
      try {
        const res = await getPharmacistDetail(row.email)
        setDetail(res.data)
      } catch {
        showAlert('error', 'Failed to load pharmacist details')
        setModal(null)
      } finally { setDetailLoading(false) }
    }
  const columns = [
    { key: 'name',  label: 'Name' },
    { key: 'email', label: 'Email' },
    { key: 'p_no',  label: 'Phone' },
    { key: 'actions', label: 'Actions', render: (row) => (
      <div className="flex gap-2">
        <button onClick={() => openDetail(row)} className="text-green-600 hover:underline text-sm">View</button>
        <button onClick={() => openEdit(row)}   className="text-blue-600 hover:underline text-sm">Edit</button>
        <button onClick={() => openDelete(row)} className="text-red-600 hover:underline text-sm">Delete</button>
      </div>
    )},
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Pharmacists"
        subtitle={`${data.length} registered pharmacists`}
        action={<button onClick={openAdd} className="btn-primary">+ Add Pharmacist</button>}
      />
      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />
      <Table columns={columns} data={data} emptyMessage="No pharmacists found" />

      {(modal === 'add' || modal === 'edit') && (
        <Modal title={modal === 'add' ? 'Add Pharmacist' : 'Edit Pharmacist'} onClose={closeModal}>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
              <input className="input-field" value={form.name} onChange={e => setForm({...form, name: e.target.value})} placeholder="Sara Malik" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input className="input-field" value={form.email} onChange={e => setForm({...form, email: e.target.value})} disabled={modal === 'edit'} placeholder="sara@hospital.com" />
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
            <div className="flex gap-3 pt-2">
              <button onClick={handleSave} disabled={saving} className="btn-primary flex-1">{saving ? 'Saving...' : 'Save'}</button>
              <button onClick={closeModal} className="btn-secondary flex-1">Cancel</button>
            </div>
          </div>
        </Modal>
      )}

      {modal === 'delete' && (
        <Modal title="Delete Pharmacist" onClose={closeModal}>
          <p className="text-gray-600 mb-6">Are you sure you want to delete <strong>{selected?.name}</strong>?</p>
          <div className="flex gap-3">
            <button onClick={handleDelete} disabled={saving} className="btn-danger flex-1">{saving ? 'Deleting...' : 'Delete'}</button>
            <button onClick={closeModal} className="btn-secondary flex-1">Cancel</button>
          </div>
        </Modal>
      )}

      {modal === 'detail' && selected && (
        <Modal title={`${selected.name} — Profile & Activity`} onClose={closeModal}>
          {detailLoading ? (
            <div className="text-center py-8 text-gray-500">Loading...</div>
          ) : detail ? (
            <div className="space-y-5">
              <div className="bg-gray-50 rounded-lg p-4 grid grid-cols-2 gap-3 text-sm">
                <div><span className="text-gray-500">Name:</span> <strong>{detail.profile?.name}</strong></div>
                <div><span className="text-gray-500">Email:</span> <strong>{detail.profile?.email}</strong></div>
                <div><span className="text-gray-500">Phone:</span> <strong>{detail.profile?.p_no}</strong></div>
              </div>
              <div>
                <h3 className="font-semibold text-gray-800 mb-2">
                  Dispensing History <span className="text-gray-400 font-normal text-sm">({detail.dispensing?.length || 0})</span>
                </h3>
                {detail.dispensing?.length === 0 ? (
                  <p className="text-gray-400 text-sm">No dispensing records found</p>
                ) : (
                  <div className="space-y-2 max-h-64 overflow-y-auto">
                    {detail.dispensing?.map((d, i) => (
                      <div key={i} className="flex justify-between items-center p-3 bg-gray-50 rounded-lg text-sm">
                        <div>
                          <span className="font-medium">{d.patient_name}</span>
                          <span className="text-gray-500 ml-2">{d.medicine_name}</span>
                        </div>
                        <div className="text-right">
                          <span className="badge-info">x{d.quantity}</span>
                          <span className="text-gray-400 ml-2">{d.date}</span>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ) : null}
        </Modal>
      )}
    </div>
  )
}
