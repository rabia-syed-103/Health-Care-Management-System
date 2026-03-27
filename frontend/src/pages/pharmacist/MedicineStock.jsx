import { useState, useEffect } from 'react'
import { getMedicineStock, getMedicineByName } from '../../api/medicines'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import LoadingSpinner from '../../components/LoadingSpinner'
import Alert from '../../components/Alert'

export default function PharmacistStock() {
  const [medicines, setMedicines] = useState([])
  const [filtered,  setFiltered]  = useState([])
  const [loading,   setLoading]   = useState(true)
  const [search,    setSearch]    = useState('')
  const [filterStatus, setFilterStatus] = useState('')
  const [alert,     setAlert]     = useState({ type: '', message: '' })

  const fetchData = async () => {
    try {
      const res = await getMedicineStock()
      const data = res.data.medicines || []
      setMedicines(data)
      setFiltered(data)
    } catch { setMedicines([]); setFiltered([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  useEffect(() => {
    let result = medicines
    if (search)       result = result.filter(m => m.name?.toLowerCase().includes(search.toLowerCase()) || m.batch_no?.toLowerCase().includes(search.toLowerCase()))
    if (filterStatus) result = result.filter(m => m.alert === filterStatus)
    setFiltered(result)
  }, [search, filterStatus, medicines])

  const handleNameSearch = async () => {
    if (!search) return fetchData()
    try {
      const res = await getMedicineByName(search)
      setFiltered(res.data.medicines || [])
    } catch {
      setFiltered([])
    }
  }

  const statusBadge = (status) => {
    const map = {
      'OUT OF STOCK':  'badge-danger',
      'EXPIRED':       'badge-danger',
      'CRITICAL':      'badge-warning',
      'EXPIRING SOON': 'badge-warning',
      'OK':            'badge-success',
    }
    return <span className={map[status] || 'badge-info'}>{status}</span>
  }

  const columns = [
    { key: 'batch_no',    label: 'Batch No' },
    { key: 'name',        label: 'Medicine Name' },
    { key: 'stock',       label: 'Stock' },
    { key: 'expiry_date', label: 'Expiry Date' },
    { key: 'alert',      label: 'Status', render: (row) => statusBadge(row.alert) },
  ]

  const statuses = ['OUT OF STOCK', 'EXPIRED', 'CRITICAL', 'EXPIRING SOON', 'OK']

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Medicine Stock"
        subtitle={`${filtered.length} of ${medicines.length} medicines`}
      />
      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

      <div className="flex gap-3 mb-4 flex-wrap">
        <div className="flex gap-2 flex-1 min-w-0">
          <input
            className="input-field flex-1"
            placeholder="Search by name or batch..."
            value={search}
            onChange={e => setSearch(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleNameSearch()}
          />
          <button onClick={handleNameSearch} className="btn-secondary">Search</button>
        </div>
        <select
          className="input-field max-w-xs"
          value={filterStatus}
          onChange={e => setFilterStatus(e.target.value)}
        >
          <option value="">All Statuses</option>
          {statuses.map(s => <option key={s} value={s}>{s}</option>)}
        </select>
        {(search || filterStatus) && (
          <button onClick={() => { setSearch(''); setFilterStatus(''); fetchData() }} className="btn-secondary">
            Clear
          </button>
        )}
      </div>

      <Table columns={columns} data={filtered} emptyMessage="No medicines found" />
    </div>
  )
}
