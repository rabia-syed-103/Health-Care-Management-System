import { useState } from 'react'
import { prescribeMedicines } from '../../api/prescriptions'
import { useAuth } from '../../context/AuthContext'
import PageHeader from '../../components/PageHeader'
import Alert from '../../components/Alert'
import { exportPDF } from '../../utils/pdfExport'

export default function DoctorPrescribe() {
  const { user } = useAuth()

  const [form, setForm] = useState({
    doctor_email:  user?.email || '',
    patient_email: '',
  })
  const [medicines, setMedicines] = useState([{ medicine_name: '', quantity: 1 }])
  const [alert,   setAlert]   = useState({ type: '', message: '' })
  const [loading, setLoading] = useState(false)
  const [result,  setResult]  = useState(null)

  const showAlert = (type, message) => {
    setAlert({ type, message })
    setTimeout(() => setAlert({ type: '', message: '' }), 5000)
  }

  const addMedicineRow = () => {
    setMedicines([...medicines, { medicine_name: '', quantity: 1 }])
  }

  const removeMedicineRow = (i) => {
    setMedicines(medicines.filter((_, idx) => idx !== i))
  }

  const updateMedicine = (i, field, value) => {
    const updated = [...medicines]
    updated[i][field] = field === 'quantity' ? parseInt(value) || 1 : value
    setMedicines(updated)
  }
  const downloadPDF = () => exportPDF(
  'Prescription Details',
  ['Field', 'Value'],
    [
      ['Prescription ID', `#${result.prescription_id}`],
      ['Patient',         result.patient],
      ['Doctor',          result.doctor],
      ['Medicines Count', result.medicines_count],
      ['Date',            new Date().toLocaleDateString()],
    ],
    `prescription-${result.prescription_id}`
  )
  const handleSubmit = async () => {
    if (!form.patient_email) return showAlert('error', 'Patient email is required')
    if (!form.doctor_email)  return showAlert('error', 'Doctor email is required')
    if (medicines.some(m => !m.medicine_name)) return showAlert('error', 'All medicine names are required')

    setLoading(true)
    setResult(null)
    try {
      const res = await prescribeMedicines({
        doctor_email:  form.doctor_email,
        patient_email: form.patient_email,
        medicines,
      })
      setResult(res.data)
      showAlert('success', 'Prescription created successfully')
      setForm({ ...form, patient_email: '' })
      setMedicines([{ medicine_name: '', quantity: 1 }])
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Failed to create prescription')
    } finally { setLoading(false) }
  }

  
  return (
    <div>
      <PageHeader title="Prescribe Medicines" subtitle="Create a new prescription for a patient" />

      <div className="max-w-2xl">
        <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

        <div className="card space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Doctor Email</label>
            <input
              className="input-field"
              value={form.doctor_email}
              onChange={e => setForm({...form, doctor_email: e.target.value})}
              placeholder="doctor@hospital.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Patient Email</label>
            <input
              className="input-field"
              value={form.patient_email}
              onChange={e => setForm({...form, patient_email: e.target.value})}
              placeholder="patient@email.com"
            />
          </div>

          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="block text-sm font-medium text-gray-700">Medicines</label>
              <button onClick={addMedicineRow} className="text-primary-600 text-sm hover:underline">+ Add Medicine</button>
            </div>

            <div className="space-y-2">
              {medicines.map((med, i) => (
                <div key={i} className="flex gap-2 items-center">
                  <input
                    className="input-field flex-1"
                    placeholder="Medicine name (e.g. Paracetamol)"
                    value={med.medicine_name}
                    onChange={e => updateMedicine(i, 'medicine_name', e.target.value)}
                  />
                  <input
                    type="number"
                    min="1"
                    className="input-field w-24"
                    placeholder="Qty"
                    value={med.quantity}
                    onChange={e => updateMedicine(i, 'quantity', e.target.value)}
                  />
                  {medicines.length > 1 && (
                    <button onClick={() => removeMedicineRow(i)} className="text-red-500 hover:text-red-700 text-xl leading-none">×</button>
                  )}
                </div>
              ))}
            </div>
          </div>

          <button onClick={handleSubmit} disabled={loading} className="btn-primary w-full">
            {loading ? 'Creating Prescription...' : 'Create Prescription'}
          </button>
        </div>

        {result && (
          <div className="card mt-4 bg-green-50 border-green-200">
            <h3 className="font-semibold text-green-800 mb-2">✅ Prescription Created</h3>
            <p className="text-sm text-green-700">Prescription ID: <strong>#{result.prescription_id}</strong></p>
            <p className="text-sm text-green-700">Patient: <strong>{result.patient}</strong></p>
            <p className="text-sm text-green-700">Doctor: <strong>{result.doctor}</strong></p>
            <p className="text-sm text-green-700">Medicines: <strong>{result.medicines_count}</strong></p>
            <button onClick={downloadPDF} className="btn-secondary mt-3">⬇ Download Prescription PDF</button>
          </div>
        )}
      </div>
    </div>
  )
}
