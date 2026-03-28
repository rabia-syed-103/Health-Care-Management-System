import { useState, useEffect } from 'react'
import { getMedicineStock, addMedicine } from '../../api/medicines'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'

const empty = { batch_no: '', name: '', stock: '', expiry_date: '' }

export default function AdminMedicines() {
  const [medicines, setMedicines] = useState([])
  const [loading,   setLoading]   = useState(true)
  const [modal,     setModal]     = useState(false)
  const [form,      setForm]      = useState(empty)
  const [alert,     setAlert]     = useState({ type: '', message: '' })
  const [saving,    setSaving]    = useState(false)
  const [search,    setSearch]    = useState('')

  const fetchData = async () => {
    try {
      const res = await getMedicineStock()
      setMedicines(res.data.medicines || [])
    } catch { setMedicines([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  const showAlert = (type, message) => { setAlert({ type, message }); setTimeout(() => setAlert({ type: '', message: '' }), 4000) }

  const handleSave = async () => {
    if (!form.batch_no || !form.name || !form.stock || !form.expiry_date) return showAlert('error', 'All fields are required')
    setSaving(true)
    try {
      await addMedicine({ ...form, stock: parseInt(form.stock) })
      showAlert('success', 'Medicine saved successfully')
      setModal(false); setForm(empty); fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Operation failed')
    } finally { setSaving(false) }
  }

  const statusBadge = (status) => {
    const map = {
      'OUT OF STOCK':   'badge-danger',
      'EXPIRED':        'badge-danger',
      'CRITICAL':       'badge-warning',
      'EXPIRING SOON':  'badge-warning',
      'OK':             'badge-success',
    }
    return <span className={map[status] || 'badge-info'}>{status}</span>
  }

  const filtered = medicines.filter(m =>
    m.name?.toLowerCase().includes(search.toLowerCase()) ||
    m.batch_no?.toLowerCase().includes(search.toLowerCase())
  )

  const columns = [
    { key: 'batch_no',    label: 'Batch No' },
    { key: 'name',        label: 'Medicine Name' },
    { key: 'stock',       label: 'Stock' },
    { key: 'expiry_date', label: 'Expiry Date' },
    { key: 'alert',      label: 'Status', render: (row) => statusBadge(row.alert) },
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Medicine Stock"
        subtitle="View and manage medicine inventory"
        action={<button onClick={() => { setForm(empty); setModal(true) }} className="btn-primary">+ Add / Restock</button>}
      />

      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

      <div className="mb-4">
        <input
          className="input-field max-w-sm"
          placeholder="Search by name or batch number..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
      </div>

      <Table columns={columns} data={filtered} emptyMessage="No medicines found" />

      {modal && (
        <Modal title="Add / Restock Medicine" onClose={() => setModal(false)}>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Batch Number</label>
              <input className="input-field" value={form.batch_no} onChange={e => setForm({...form, batch_no: e.target.value})} placeholder="B001" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Medicine Name</label>
              <input className="input-field" value={form.name} onChange={e => setForm({...form, name: e.target.value})} placeholder="Paracetamol" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Stock Quantity</label>
              <input type="number" className="input-field" value={form.stock} onChange={e => setForm({...form, stock: e.target.value})} placeholder="100" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Expiry Date</label>
              <input type="date" className="input-field" value={form.expiry_date} onChange={e => setForm({...form, expiry_date: e.target.value})} />
            </div>
            <p className="text-xs text-gray-400">If batch number already exists, stock will be added on top.</p>
            <div className="flex gap-3 pt-2">
              <button onClick={handleSave} disabled={saving} className="btn-primary flex-1">{saving ? 'Saving...' : 'Save'}</button>
              <button onClick={() => setModal(false)} className="btn-secondary flex-1">Cancel</button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  )
}
