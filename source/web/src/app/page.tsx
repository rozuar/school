'use client'

import { useEffect, useState } from 'react'
import axios from 'axios'
import { PersonIcon } from '../components/PersonIcon'

const API_URL = '/api/v1'

interface Usuario {
  id: string
  email: string
  nombre: string
  rol: string
}

interface Curso {
  id: string
  nombre: string
  nivel: string
}

interface Alumno {
  id: string
  curso_id: string
  nombre: string
  apellido: string
  rut: string
  caso_especial: boolean
  activo: boolean
  estado_temporal?: string // bano, enfermeria, sos
}

interface BloqueHorario {
  id: string
  numero: number
  hora_inicio: string
  hora_fin: string
}

interface Horario {
  id: string
  curso_id: string
  asignatura_id: string
  profesor_id: string
  bloque_id: string
  dia_semana: number
  asignatura?: { id: string; nombre: string }
  curso?: { id: string; nombre: string; nivel: string }
  bloque?: BloqueHorario
}

interface EstadoTemporal {
  id: string
  alumno_id: string
  tipo: string
  inicio: string
  fin?: string
  alumno?: Alumno
}

export default function Home() {
  const [token, setToken] = useState<string | null>(null)
  const [usuario, setUsuario] = useState<Usuario | null>(null)
  const [alumnos, setAlumnos] = useState<Alumno[]>([])
  const [horarios, setHorarios] = useState<Horario[]>([])
  const [horarioSeleccionado, setHorarioSeleccionado] = useState<Horario | null>(null)
  const [vista, setVista] = useState<'horario' | 'clase'>('horario')
  const [estadosTemporales, setEstadosTemporales] = useState<EstadoTemporal[]>([])
  const [diaSeleccionado, setDiaSeleccionado] = useState<number>(() => {
    const jsDay = new Date().getDay() // 0=domingo..6=sabado
    if (jsDay >= 1 && jsDay <= 5) return jsDay // lunes..viernes
    return 1
  })

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(true)
  const [loginLoading, setLoginLoading] = useState(false)
  const [expandedAlumnoId, setExpandedAlumnoId] = useState<string | null>(null)

  // Asistencia del bloque actual
  const [asistencia, setAsistencia] = useState<Record<string, string>>({}) // alumno_id -> estado

  const getAlumnoColor = (alumno: Alumno) => {
    const estadoTemp = estadosTemporales.find(e => e.alumno_id === alumno.id && !e.fin)
    if (estadoTemp) {
      if (estadoTemp.tipo === 'sos') return '#dc2626' // rojo
      return '#ca8a04' // amarillo para bano/enfermeria
    }
    const estado = asistencia[alumno.id]
    if (estado === 'ausente') return '#111827' // negro
    return '#16a34a' // verde
  }

  const getEstadoTexto = (alumno: Alumno) => {
    const estadoTemp = estadosTemporales.find(e => e.alumno_id === alumno.id && !e.fin)
    if (estadoTemp) {
      if (estadoTemp.tipo === 'bano') return '🚻 Baño'
      if (estadoTemp.tipo === 'enfermeria') return '🏥 Enfermería'
      if (estadoTemp.tipo === 'sos') return '🆘 SOS'
    }
    const estado = asistencia[alumno.id]
    if (estado === 'ausente') return 'Ausente'
    if (estado === 'presente') return 'Presente'
    return 'Sin registrar'
  }

  // Cargar token al inicio
  useEffect(() => {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      setToken(storedToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
      cargarUsuario()
    } else {
      setLoading(false)
    }
  }, [])

  const cargarUsuario = async () => {
    try {
      const resp = await axios.get(`${API_URL}/auth/me`)
      setUsuario(resp.data)
      await Promise.all([cargarHorariosMis(), cargarEstadosTemporales()])
    } catch (error) {
      console.error('Error cargando usuario:', error)
      logout()
    } finally {
      setLoading(false)
    }
  }

  const cargarHorariosMis = async () => {
    try {
      const resp = await axios.get(`${API_URL}/horarios/mis`)
      setHorarios(resp.data || [])
    } catch (error) {
      console.error('Error cargando horarios del profesor:', error)
    }
  }

  const cargarAlumnos = async (cursoId: string) => {
    try {
      const resp = await axios.get(`${API_URL}/cursos/${cursoId}/alumnos`)
      setAlumnos(resp.data || [])
      // Inicializar asistencia como "presente" por defecto
      const asistenciaInicial: Record<string, string> = {}
      ;(resp.data || []).forEach((a: Alumno) => {
        asistenciaInicial[a.id] = 'presente'
      })
      setAsistencia(asistenciaInicial)
    } catch (error) {
      console.error('Error cargando alumnos:', error)
    }
  }

  const cargarEstadosTemporales = async () => {
    try {
      const resp = await axios.get(`${API_URL}/estados-temporales`)
      setEstadosTemporales(resp.data || [])
    } catch (error) {
      console.error('Error cargando estados temporales:', error)
    }
  }

  const login = async () => {
    if (!email || !password) {
      alert('Ingresa email y contraseña')
      return
    }
    setLoginLoading(true)
    try {
      const resp = await axios.post(`${API_URL}/auth/login`, { email, password })
      const { token: newToken, usuario: user } = resp.data
      localStorage.setItem('token', newToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      setToken(newToken)
      setUsuario(user)
      await Promise.all([cargarHorariosMis(), cargarEstadosTemporales()])
    } catch (error: any) {
      console.error('Error en login:', error)
      alert(error?.response?.data?.error || 'Error en login')
    } finally {
      setLoginLoading(false)
    }
  }

  const loginDemo = async () => {
    setLoginLoading(true)
    try {
      // Primero hacer seed si no hay datos
      try {
        await axios.post(`${API_URL}/seed`)
      } catch (e) {
        // Ignorar error si ya hay datos
      }
      // Login con usuario demo
      const resp = await axios.post(`${API_URL}/auth/login`, {
        email: 'profesor1@escuela.cl',
        password: 'profesor123'
      })
      const { token: newToken, usuario: user } = resp.data
      localStorage.setItem('token', newToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      setToken(newToken)
      setUsuario(user)
      await Promise.all([cargarHorariosMis(), cargarEstadosTemporales()])
    } catch (error: any) {
      console.error('Error en login demo:', error)
      alert(error?.response?.data?.error || 'Error en login demo')
    } finally {
      setLoginLoading(false)
    }
  }

  const logout = () => {
    localStorage.removeItem('token')
    delete axios.defaults.headers.common['Authorization']
    setToken(null)
    setUsuario(null)
    setAlumnos([])
    setHorarios([])
    setHorarioSeleccionado(null)
    setVista('horario')
  }

  const cargarAsistenciaHorarioHoy = async (horarioId: string) => {
    try {
      const fecha = new Date().toISOString().split('T')[0]
      const resp = await axios.get(`${API_URL}/asistencia/horario/${horarioId}/fecha/${fecha}`)
      const rows = resp.data || []
      if (rows.length === 0) {
        // si no hay registros guardados para este bloque hoy, dejamos default "presente"
        return
      }
      const map: Record<string, string> = {}
      rows.forEach((r: any) => {
        if (r?.alumno_id && r?.estado) map[r.alumno_id] = r.estado
      })
      setAsistencia(prev => ({ ...prev, ...map }))
    } catch (error) {
      // si no existe aún, no hacemos nada
    }
  }

  const toggleAsistencia = (alumnoId: string) => {
    setAsistencia(prev => ({
      ...prev,
      [alumnoId]: prev[alumnoId] === 'presente' ? 'ausente' : 'presente'
    }))
  }

  const guardarAsistencia = async () => {
    if (!horarioSeleccionado) {
      alert('Selecciona un bloque horario primero')
      return
    }

    const registros = Object.entries(asistencia).map(([alumno_id, estado]) => ({
      alumno_id,
      estado
    }))

    try {
      await axios.post(`${API_URL}/asistencia/bloque`, {
        horario_id: horarioSeleccionado.id,
        fecha: new Date().toISOString().split('T')[0],
        registros
      })
      alert('Asistencia guardada correctamente')
    } catch (error: any) {
      console.error('Error guardando asistencia:', error)
      alert(error?.response?.data?.error || 'Error guardando asistencia')
    }
  }

  const setEstadoTemporal = async (alumnoId: string, tipo: string) => {
    try {
      await axios.put(`${API_URL}/alumnos/${alumnoId}/estado-temporal`, { tipo })
      await cargarEstadosTemporales()
    } catch (error: any) {
      console.error('Error estableciendo estado temporal:', error)
      alert(error?.response?.data?.error || 'Error estableciendo estado')
    }
  }

  const clearEstadoTemporal = async (alumnoId: string) => {
    try {
      await axios.delete(`${API_URL}/alumnos/${alumnoId}/estado-temporal`)
      await cargarEstadosTemporales()
    } catch (error: any) {
      console.error('Error limpiando estado temporal:', error)
      alert(error?.response?.data?.error || 'Error limpiando estado')
    }
  }

  const minutos = (hhmm?: string) => {
    if (!hhmm) return null
    const [h, m] = hhmm.split(':').map(Number)
    if (Number.isNaN(h) || Number.isNaN(m)) return null
    return h * 60 + m
  }

  const esBloqueActual = (h: Horario) => {
    const hoy = new Date()
    const jsDay = hoy.getDay()
    if (jsDay !== h.dia_semana) return false
    const nowMin = hoy.getHours() * 60 + hoy.getMinutes()
    const ini = minutos(h.bloque?.hora_inicio)
    const fin = minutos(h.bloque?.hora_fin)
    if (ini == null || fin == null) return false
    return nowMin >= ini && nowMin < fin
  }

  const diasSemana = [
    { id: 1, label: 'Lunes' },
    { id: 2, label: 'Martes' },
    { id: 3, label: 'Miércoles' },
    { id: 4, label: 'Jueves' },
    { id: 5, label: 'Viernes' },
  ]

  const horariosDia = horarios
    .filter((h) => h.dia_semana === diaSeleccionado)
    .slice()
    .sort((a, b) => (a.bloque?.numero || 0) - (b.bloque?.numero || 0))

  const diaLabel = diasSemana.find(d => d.id === diaSeleccionado)?.label || ''

  if (loading) {
    return (
      <div style={{ padding: '2rem', textAlign: 'center' }}>
        <p>Cargando...</p>
      </div>
    )
  }

  // Pantalla de login
  if (!token || !usuario) {
    return (
      <div style={{ padding: '2rem', textAlign: 'center', maxWidth: 400, margin: '0 auto' }}>
        <h1 style={{ marginBottom: '1rem' }}>Plataforma Escolar</h1>
        <p style={{ color: '#6b7280', marginBottom: '1.5rem' }}>
          Ingresa con tus credenciales
        </p>

        <div style={{ display: 'grid', gap: '0.75rem' }}>
          <button
            onClick={loginDemo}
            disabled={loginLoading}
            style={{
              padding: '0.75rem 1.5rem',
              backgroundColor: '#2563eb',
              color: 'white',
              border: 'none',
              borderRadius: 8,
              cursor: 'pointer',
              fontWeight: 700,
            }}
          >
            {loginLoading ? 'Ingresando...' : 'Login Demo (Profesor)'}
          </button>

          <div style={{ color: '#9ca3af', fontSize: '0.875rem', margin: '0.5rem 0' }}>
            o ingresa manualmente:
          </div>

          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Email"
            style={{
              padding: '0.75rem 1rem',
              borderRadius: 8,
              border: '1px solid #e5e7eb',
            }}
          />
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Contraseña"
            style={{
              padding: '0.75rem 1rem',
              borderRadius: 8,
              border: '1px solid #e5e7eb',
            }}
          />
          <button
            onClick={login}
            disabled={loginLoading}
            style={{
              padding: '0.75rem 1.5rem',
              backgroundColor: '#111827',
              color: 'white',
              border: 'none',
              borderRadius: 8,
              cursor: 'pointer',
              fontWeight: 600,
            }}
          >
            Ingresar
          </button>
        </div>
      </div>
    )
  }

  // Pantalla principal (horario del profesor)
  return (
    <main style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
      {vista === 'horario' ? (
        <>
          <header style={{ marginBottom: '2rem', display: 'flex', justifyContent: 'space-between' }}>
            <div>
              <h1>{usuario.nombre}</h1>
              <p style={{ color: '#6b7280' }}>Selecciona un bloque para ingresar a la clase.</p>
            </div>
            <button
              onClick={logout}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#ef4444',
                color: 'white',
                border: 'none',
                borderRadius: 8,
                cursor: 'pointer',
                height: 'fit-content',
              }}
            >
              Cerrar sesion
            </button>
          </header>

          {/* Horario semanal (L-V) */}
          <section style={{ marginBottom: '2rem' }}>
            <h2>Horario del profesor</h2>

            <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', marginTop: '0.75rem' }}>
              {diasSemana.map((d) => (
                <button
                  key={d.id}
                  onClick={() => setDiaSeleccionado(d.id)}
                  style={{
                    padding: '0.5rem 0.9rem',
                    borderRadius: 999,
                    border: '1px solid #e5e7eb',
                    background: diaSeleccionado === d.id ? '#111827' : '#fff',
                    color: diaSeleccionado === d.id ? '#fff' : '#111827',
                    cursor: 'pointer',
                    fontWeight: 800,
                  }}
                >
                  {d.label}
                </button>
              ))}
            </div>

            <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', marginTop: '0.75rem' }}>
              {horariosDia.map((h) => {
                const now = esBloqueActual(h)
                return (
                  <button
                    key={h.id}
                    onClick={async () => {
                      setHorarioSeleccionado(h)
                      setVista('clase')
                      await Promise.all([
                        cargarAlumnos(h.curso_id),
                        cargarAsistenciaHorarioHoy(h.id),
                        cargarEstadosTemporales(),
                      ])
                    }}
                    style={{
                      padding: '0.75rem 1rem',
                      backgroundColor: 'white',
                      color: '#111827',
                      border: now ? '2px solid #f59e0b' : '1px solid #e5e7eb',
                      borderRadius: 10,
                      cursor: 'pointer',
                      minWidth: 240,
                      textAlign: 'left',
                    }}
                  >
                    <div style={{ fontWeight: 800 }}>
                      {diaLabel} · Bloque {h.bloque?.numero} {now ? '• Ahora' : ''}
                    </div>
                    <div style={{ fontSize: '0.85rem' }}>
                      {h.bloque?.hora_inicio} - {h.bloque?.hora_fin}
                    </div>
                    <div style={{ fontSize: '0.85rem', marginTop: '0.25rem' }}>
                      {h.asignatura?.nombre || '—'} · {h.curso?.nombre || '—'}
                    </div>
                  </button>
                )
              })}

              {horarios.length === 0 && (
                <p style={{ color: '#6b7280' }}>No hay horarios asignados a este profesor.</p>
              )}
              {horarios.length > 0 && horariosDia.length === 0 && (
                <p style={{ color: '#6b7280' }}>No hay bloques para este día.</p>
              )}
            </div>
          </section>
        </>
      ) : (
        <>
          <header style={{ marginBottom: '1.25rem', display: 'flex', justifyContent: 'space-between', alignItems: 'start', gap: '1rem' }}>
            <div style={{ minWidth: 0 }}>
              <button
                onClick={() => setVista('horario')}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: '#f3f4f6',
                  border: 'none',
                  borderRadius: 8,
                  cursor: 'pointer',
                  marginBottom: '0.75rem',
                  fontWeight: 800,
                }}
              >
                ← Volver al horario
              </button>
              <h1 style={{ marginBottom: '0.25rem' }}>{usuario.nombre}</h1>
              <div style={{ color: '#6b7280', fontSize: '0.95rem' }}>
                {horarioSeleccionado
                  ? `${diaLabel} · ${horarioSeleccionado.bloque?.hora_inicio}-${horarioSeleccionado.bloque?.hora_fin} · ${horarioSeleccionado.asignatura?.nombre || '—'} · ${horarioSeleccionado.curso?.nombre || '—'}`
                  : 'Clase'}
              </div>
            </div>
            <button
              onClick={logout}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#ef4444',
                color: 'white',
                border: 'none',
                borderRadius: 8,
                cursor: 'pointer',
                height: 'fit-content',
              }}
            >
              Cerrar sesion
            </button>
          </header>

          {/* Lista de alumnos (clase) */}
          <section>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2>Alumnos ({alumnos.length})</h2>
          <button
            onClick={guardarAsistencia}
            disabled={!horarioSeleccionado}
            style={{
              padding: '0.75rem 1.5rem',
              backgroundColor: horarioSeleccionado ? '#22c55e' : '#9ca3af',
              color: 'white',
              border: 'none',
              borderRadius: 8,
              cursor: horarioSeleccionado ? 'pointer' : 'not-allowed',
              fontWeight: 700,
            }}
          >
            Guardar Asistencia
          </button>
        </div>

        {!horarioSeleccionado ? (
          <div style={{ marginTop: '1rem', padding: '1rem', background: '#fff', borderRadius: 8 }}>
            Selecciona un bloque del horario para ver alumnos y registrar asistencia.
          </div>
        ) : alumnos.length === 0 ? (
          <div style={{ marginTop: '1rem', padding: '1rem', background: '#fff', borderRadius: 8 }}>
            Sin alumnos en este curso.
          </div>
        ) : (
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
              gap: '0.75rem',
              marginTop: '1rem',
            }}
          >
            {alumnos.map((alumno) => {
              const expanded = expandedAlumnoId === alumno.id
              const tieneEstadoTemp = estadosTemporales.some(e => e.alumno_id === alumno.id && !e.fin)

              return (
                <div
                  key={alumno.id}
                  style={{
                    background: '#fff',
                    border: '1px solid #e5e7eb',
                    borderRadius: 10,
                    padding: '0.9rem',
                  }}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', gap: '0.75rem', alignItems: 'start' }}>
                    <div style={{ minWidth: 0 }}>
                      <div style={{ fontWeight: 800, display: 'flex', alignItems: 'center', gap: 8 }}>
                        <PersonIcon color={getAlumnoColor(alumno)} title={getEstadoTexto(alumno)} />
                        <span style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                          {alumno.nombre} {alumno.apellido}
                        </span>
                      </div>
                      <div style={{ color: '#6b7280', fontSize: '0.85rem', marginTop: '0.15rem' }}>
                        {getEstadoTexto(alumno)}
                        {alumno.caso_especial && ' - Caso Especial'}
                      </div>
                    </div>

                    <button
                      onClick={() => setExpandedAlumnoId(expanded ? null : alumno.id)}
                      style={{
                        padding: '0.35rem 0.6rem',
                        borderRadius: 8,
                        border: '1px solid #e5e7eb',
                        background: '#fff',
                        cursor: 'pointer',
                        fontWeight: 800,
                      }}
                    >
                      {expanded ? '▾' : '▸'}
                    </button>
                  </div>

                  {expanded && (
                    <div style={{ marginTop: '0.75rem', display: 'flex', gap: '0.35rem', flexWrap: 'wrap' }}>
                      <button
                        onClick={() => toggleAsistencia(alumno.id)}
                        style={{
                          padding: '0.35rem 0.6rem',
                          borderRadius: 8,
                          border: '1px solid #e5e7eb',
                          background: asistencia[alumno.id] === 'presente' ? '#dcfce7' : '#fef2f2',
                          cursor: 'pointer',
                          fontWeight: 800,
                        }}
                      >
                        {asistencia[alumno.id] === 'presente' ? '✅ Presente' : '❌ Ausente'}
                      </button>

                      {!tieneEstadoTemp ? (
                        <>
                          <button
                            onClick={() => setEstadoTemporal(alumno.id, 'bano')}
                            style={{
                              padding: '0.35rem 0.6rem',
                              borderRadius: 8,
                              border: '1px solid #e5e7eb',
                              background: '#fef9c3',
                              cursor: 'pointer',
                              fontWeight: 700,
                            }}
                          >
                            🚻 Baño
                          </button>
                          <button
                            onClick={() => setEstadoTemporal(alumno.id, 'enfermeria')}
                            style={{
                              padding: '0.35rem 0.6rem',
                              borderRadius: 8,
                              border: '1px solid #e5e7eb',
                              background: '#fef9c3',
                              cursor: 'pointer',
                              fontWeight: 700,
                            }}
                          >
                            🏥 Enfermería
                          </button>
                          <button
                            onClick={() => setEstadoTemporal(alumno.id, 'sos')}
                            style={{
                              padding: '0.35rem 0.6rem',
                              borderRadius: 8,
                              border: '1px solid #e5e7eb',
                              background: '#fee2e2',
                              cursor: 'pointer',
                              fontWeight: 800,
                            }}
                          >
                            🆘 SOS
                          </button>
                        </>
                      ) : (
                        <button
                          onClick={() => clearEstadoTemporal(alumno.id)}
                          style={{
                            padding: '0.35rem 0.6rem',
                            borderRadius: 8,
                            border: '1px solid #e5e7eb',
                            background: '#dbeafe',
                            cursor: 'pointer',
                            fontWeight: 700,
                          }}
                        >
                          ↩️ Volvio
                        </button>
                      )}
                    </div>
                  )}
                </div>
              )
            })}
          </div>
        )}
      </section>
        </>
      )}
    </main>
  )
}






