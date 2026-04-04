import { useState } from 'react'
import { getPatientHistory } from '../../api/doctor'
import PageHeader from '../../components/PageHeader'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'
import { exportPDF } from '../../utils/pdfExport'

export default function DoctorPatientHistory() {
  const [email,   setEmail]   = useState('')
  const [data,    setData]    = useState(null)
  const [loading, setLoading] = useState(false)
  const [alert,   setAlert]   = useState({ type: '', message: '' })

  const showAlert = (type, message) => {
    setAlert({ type, message })
    setTimeout(() => setAlert({ type: '', message: '' }), 5000)
  }

  const handleSearch = async () => {
    if (!email) return showAlert('error', 'Please enter a patient email')
    setLoading(true)
    setData(null)
    try {
      const res = await getPatientHistory(email)
      setData(res.data)
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Patient not found')
    } finally { setLoading(false) }
  }
  const downloadPDF = () => {
  const rows = []

  // Appointments section
  data.appointments?.forEach(a => {
    rows.push([a.doctor_name, a.date, a.time, a.status, '—', '—', '—'])
  })

  // Prescriptions section
  data.prescriptions?.forEach(p => {
    p.medicines?.forEach(m => {
      rows.push(['—', p.date, '—', '—', `Dr. ${p.doctor_name}`, m.medicine_name, m.quantity])
    })
  })

  exportPDF(
    `Patient History — ${data.patient?.name}`,
    ['Doctor', 'Date', 'Time', 'Appt Status', 'Prescribing Dr', 'Medicine', 'Qty'],
    rows,
    `patient-history-${data.patient?.name?.replace(' ', '-')}`
  )
}
  return (
    <div>
      <PageHeader title="Patient History" subtitle="Search full history by patient email" />

      <div className="max-w-2xl">
        <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

        <div className="flex gap-3 mb-6">
          <input
            className="input-field flex-1"
            placeholder="patient@email.com"
            value={email}
            onChange={e => setEmail(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleSearch()}
          />
          <button onClick={handleSearch} disabled={loading} className="btn-primary">
            {loading ? 'Searching...' : 'Search'}
          </button>
        </div>

        {loading && <LoadingSpinner text="Fetching patient history..." />}

        {data && (
          <div className="space-y-6">
            <div className="flex justify-end">
              <button onClick={downloadPDF} className="btn-secondary">⬇ Download PDF</button>
            </div>
            {/* Patient Info */}
            <div className="card">
              <h2 className="text-lg font-semibold text-gray-800 mb-3">Patient Info</h2>
              <div className="grid grid-cols-2 gap-3 text-sm">
                <div><span className="text-gray-500">Name:</span> <strong>{data.patient?.name}</strong></div>
                <div><span className="text-gray-500">Blood Group:</span> <strong>{data.patient?.b_gr}</strong></div>
                <div><span className="text-gray-500">Phone:</span> <strong>{data.patient?.p_no}</strong></div>
                <div><span className="text-gray-500">ID:</span> <strong>#{data.patient?.id}</strong></div>
              </div>
            </div>

            {/* Appointments */}
            <div className="card">
              <h2 className="text-lg font-semibold text-gray-800 mb-3">
                Appointments <span className="text-gray-400 text-sm font-normal">({data.appointments?.length || 0})</span>
              </h2>
              {data.appointments?.length === 0 ? (
                <p className="text-gray-400 text-sm">No appointments found</p>
              ) : (
                <div className="space-y-2">
                  {data.appointments?.map((a, i) => (
                    <div key={i} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg text-sm">
                      <div>
                        <span className="font-medium">{a.doctor_name}</span>
                        <span className="text-gray-500 ml-2">{a.date} at {a.time}</span>
                      </div>
                      <span className={a.status === 'pending' ? 'badge-warning' : 'badge-success'}>{a.status}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Prescriptions */}
            <div className="card">
              <h2 className="text-lg font-semibold text-gray-800 mb-3">
                Prescriptions <span className="text-gray-400 text-sm font-normal">({data.prescriptions?.length || 0})</span>
              </h2>
              {data.prescriptions?.length === 0 ? (
                <p className="text-gray-400 text-sm">No prescriptions found</p>
              ) : (
                <div className="space-y-3">
                  {data.prescriptions?.map((p, i) => (
                    <div key={i} className="p-3 bg-gray-50 rounded-lg text-sm">
                      <div className="flex justify-between mb-2">
                        <span className="font-medium">Dr. {p.doctor_name}</span>
                        <span className="text-gray-500">{p.date}</span>
                      </div>
                      <div className="space-y-1">
                        {p.medicines?.map((m, j) => (
                          <div key={j} className="flex justify-between text-gray-600 pl-2">
                            <span>• {m.medicine_name}</span>
                            <span>x{m.quantity}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
