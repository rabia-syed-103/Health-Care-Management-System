import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './context/AuthContext'

import Login        from './pages/auth/Login'
import Register     from './pages/auth/Register'
import Unauthorized from './pages/Unauthorized'
import ProtectedRoute from './components/ProtectedRoute'
import Layout         from './components/Layout'

// Admin pages
import AdminDashboard    from './pages/admin/Dashboard'
import AdminDoctors      from './pages/admin/Doctors'
import AdminReceptionists from './pages/admin/Receptionists'
import AdminPharmacists  from './pages/admin/Pharmacists'
import AdminBloodManagers from './pages/admin/BloodManagers'
import AdminMedicines    from './pages/admin/Medicines'

// Doctor pages
import DoctorDashboard    from './pages/doctor/Dashboard'
import DoctorPrescribe    from './pages/doctor/Prescribe'
import DoctorBloodRequest from './pages/doctor/BloodRequest'
import DoctorPatientHistory from './pages/doctor/PatientHistory'

// Receptionist pages
import ReceptionistDashboard    from './pages/receptionist/Dashboard'
import ReceptionistPatients     from './pages/receptionist/Patients'
import ReceptionistAppointments from './pages/receptionist/Appointments'
import ReceptionistERPatients   from './pages/receptionist/ERPatients'
import ReceptionistEREntries    from './pages/receptionist/EREntries'

// Pharmacist pages
import PharmacistDashboard  from './pages/pharmacist/Dashboard'
import PharmacistPending    from './pages/pharmacist/PendingPrescriptions'
import PharmacistStock      from './pages/pharmacist/MedicineStock'
import PharmacistAddMed     from './pages/pharmacist/AddMedicine'
import PharmacistDispense   from './pages/pharmacist/Dispense'

// Blood Manager pages
import BloodManagerDashboard  from './pages/blood_manager/Dashboard'
import BloodManagerDonors     from './pages/blood_manager/Donors'
import BloodManagerDonate     from './pages/blood_manager/RecordDonation'
import BloodManagerInventory  from './pages/blood_manager/Inventory'
import BloodManagerRequests   from './pages/blood_manager/PendingRequests'
import BloodManagerExpired    from './pages/blood_manager/ExpiredBlood'

function RoleRedirect() {
  const { role } = useAuth()
  const redirects = {
    admin:         '/admin',
    doctor:        '/doctor',
    receptionist:  '/receptionist',
    pharmacist:    '/pharmacist',
    blood_manager: '/blood-manager',
  }
  return <Navigate to={redirects[role] || '/login'} replace />
}

function WrappedRoute({ component: Component, roles }) {
  return (
    <ProtectedRoute allowedRoles={roles}>
      <Layout>
        <Component />
      </Layout>
    </ProtectedRoute>
  )
}

export default function App() {
  return (
    <Routes>
      {/* Public */}
      <Route path="/login"        element={<Login />} />
      <Route path="/register"     element={<Register />} />
      <Route path="/unauthorized" element={<Unauthorized />} />
      <Route path="/"             element={<ProtectedRoute><RoleRedirect /></ProtectedRoute>} />

      {/* Admin */}
      <Route path="/admin"               element={<WrappedRoute component={AdminDashboard}     roles={['admin']} />} />
      <Route path="/admin/doctors"       element={<WrappedRoute component={AdminDoctors}       roles={['admin']} />} />
      <Route path="/admin/receptionists" element={<WrappedRoute component={AdminReceptionists} roles={['admin']} />} />
      <Route path="/admin/pharmacists"   element={<WrappedRoute component={AdminPharmacists}   roles={['admin']} />} />
      <Route path="/admin/blood-managers"element={<WrappedRoute component={AdminBloodManagers} roles={['admin']} />} />
      <Route path="/admin/medicines"     element={<WrappedRoute component={AdminMedicines}     roles={['admin']} />} />

      {/* Doctor */}
      <Route path="/doctor"                element={<WrappedRoute component={DoctorDashboard}     roles={['doctor','admin']} />} />
      <Route path="/doctor/prescribe"      element={<WrappedRoute component={DoctorPrescribe}     roles={['doctor','admin']} />} />
      <Route path="/doctor/blood-request"  element={<WrappedRoute component={DoctorBloodRequest}  roles={['doctor','admin']} />} />
      <Route path="/doctor/patient-history"element={<WrappedRoute component={DoctorPatientHistory}roles={['doctor','admin']} />} />

      {/* Receptionist */}
      <Route path="/receptionist"              element={<WrappedRoute component={ReceptionistDashboard}    roles={['receptionist','admin']} />} />
      <Route path="/receptionist/patients"     element={<WrappedRoute component={ReceptionistPatients}     roles={['receptionist','admin']} />} />
      <Route path="/receptionist/appointments" element={<WrappedRoute component={ReceptionistAppointments} roles={['receptionist','admin']} />} />
      <Route path="/receptionist/er-patients"  element={<WrappedRoute component={ReceptionistERPatients}   roles={['receptionist','admin']} />} />
      <Route path="/receptionist/er-entries"   element={<WrappedRoute component={ReceptionistEREntries}    roles={['receptionist','admin']} />} />

      {/* Pharmacist */}
      <Route path="/pharmacist"              element={<WrappedRoute component={PharmacistDashboard} roles={['pharmacist','admin']} />} />
      <Route path="/pharmacist/pending"      element={<WrappedRoute component={PharmacistPending}   roles={['pharmacist','admin']} />} />
      <Route path="/pharmacist/stock"        element={<WrappedRoute component={PharmacistStock}     roles={['pharmacist','admin']} />} />
      <Route path="/pharmacist/add-medicine" element={<WrappedRoute component={PharmacistAddMed}    roles={['pharmacist','admin']} />} />
      <Route path="/pharmacist/dispense"     element={<WrappedRoute component={PharmacistDispense}  roles={['pharmacist','admin']} />} />

      {/* Blood Manager */}
      <Route path="/blood-manager"           element={<WrappedRoute component={BloodManagerDashboard} roles={['blood_manager','admin']} />} />
      <Route path="/blood-manager/donors"    element={<WrappedRoute component={BloodManagerDonors}    roles={['blood_manager','admin']} />} />
      <Route path="/blood-manager/donate"    element={<WrappedRoute component={BloodManagerDonate}    roles={['blood_manager','admin']} />} />
      <Route path="/blood-manager/inventory" element={<WrappedRoute component={BloodManagerInventory} roles={['blood_manager','admin']} />} />
      <Route path="/blood-manager/requests"  element={<WrappedRoute component={BloodManagerRequests}  roles={['blood_manager','admin']} />} />
      <Route path="/blood-manager/expired"   element={<WrappedRoute component={BloodManagerExpired}   roles={['blood_manager','admin']} />} />
    </Routes>
  )
}