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
  activo?: boolean
}

interface Concepto {
  id: string
  codigo: string
  nombre: string
  descripcion: string
  activo: boolean
}

interface Accion {
  id: string
  codigo: string
  nombre: string
  tipo: string
  parametros: any
  activo: boolean
}

interface Regla {
  id: string
  nombre: string
  concepto_id: string
  concepto?: Concepto
  condicion: any
  accion_id: string
  accion?: Accion
  activo: boolean
}

interface Evento {
  id: string
  concepto_id: string
  concepto?: Concepto
  alumno_id?: string
  alumno?: { nombre: string; apellido: string }
  curso_id?: string
  curso?: { nombre: string }
  origen: string
  activo: boolean
  created_at: string
}

interface DashboardData {
  total_cursos: number
  total_alumnos: number
  eventos_activos: number
  estados_temporales_activos: number
  asistencia_hoy: {
    presentes: number
    ausentes: number
    justificados: number
  }
  ultimos_eventos: Evento[]
}

interface Alumno {
  id: string
  curso_id: string
  nombre: string
  apellido: string
  caso_especial: boolean
  activo: boolean
}

interface EstadoTemporal {
  id: string
  alumno_id: string
  tipo: string
  inicio: string
  fin?: string
  alumno?: Alumno
}

interface Curso {
  id: string
  nombre: string
  nivel: string
}

interface MonitorCurso {
  curso_id: string
  nombre: string
  nivel: string
  sala_semaforo: 'verde' | 'amarillo' | 'rojo' | 'gris'
  profesor_semaforo: 'verde' | 'amarillo' | 'gris'
  eventos_activos: number
  estados_temporales_activos: number
  ultima_asistencia_en?: string
}

interface Alerta {
  id: string
  codigo: string
  titulo: string
  prioridad: string
  estado: string
  curso_id?: string
  alumno_id?: string
  evento_id?: string
  created_at: string
}

interface Asignatura {
  id: string
  nombre: string
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
  dia_semana: number
  bloque_id: string
  asignatura_id: string
  profesor_id: string
  asignatura?: Asignatura
  profesor?: Usuario
  bloque?: BloqueHorario
}

interface Auditoria {
  id: string
  tabla: string
  registro_id: string
  accion: string
  datos_anteriores?: any
  datos_nuevos?: any
  usuario_id?: string
  usuario?: Usuario
  created_at: string
}

interface AccionEjecucion {
  id: string
  regla_id: string
  accion_id: string
  evento_id: string
  alumno_id?: string
  curso_id?: string
  resultado: string
  detalle?: any
  ejecutado_en: string
  regla?: Regla
  accion?: Accion
  created_at: string
}

type Tab =
  | 'dashboard'
  | 'monitor'
  | 'profesores'
  | 'horarios'
  | 'trazabilidad'
  | 'rbac'
  | 'conceptos'
  | 'acciones'
  | 'reglas'
  | 'eventos'

export default function Backoffice() {
  const [token, setToken] = useState<string | null>(null)
  const [usuario, setUsuario] = useState<Usuario | null>(null)
  const [loading, setLoading] = useState(true)
  const [loginLoading, setLoginLoading] = useState(false)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const [tab, setTab] = useState<Tab>('dashboard')
  const [dashboard, setDashboard] = useState<DashboardData | null>(null)
  const [conceptos, setConceptos] = useState<Concepto[]>([])
  const [acciones, setAcciones] = useState<Accion[]>([])
  const [reglas, setReglas] = useState<Regla[]>([])
  const [eventos, setEventos] = useState<Evento[]>([])

  // Profesores
  const [profesores, setProfesores] = useState<Usuario[]>([])
  const [nuevoProfesorNombre, setNuevoProfesorNombre] = useState('')
  const [nuevoProfesorEmail, setNuevoProfesorEmail] = useState('')
  const [nuevoProfesorPassword, setNuevoProfesorPassword] = useState('profesor123')

  // Horarios
  const [cursos, setCursos] = useState<Curso[]>([])
  const [asignaturas, setAsignaturas] = useState<Asignatura[]>([])
  const [bloques, setBloques] = useState<BloqueHorario[]>([])
  const [cursoHorarioId, setCursoHorarioId] = useState<string>('')
  const [horarios, setHorarios] = useState<Horario[]>([])
  const [savingHorarioKey, setSavingHorarioKey] = useState<string | null>(null)
  const [csvImport, setCsvImport] = useState<string>('curso,dia_semana,bloque_numero,asignatura,profesor_email\n')
  const [csvImportLoading, setCsvImportLoading] = useState(false)

  // Trazabilidad
  const [auditorias, setAuditorias] = useState<Auditoria[]>([])
  const [accionesEjecuciones, setAccionesEjecuciones] = useState<AccionEjecucion[]>([])
  const [auditTabla, setAuditTabla] = useState<string>('')
  const [traceLoading, setTraceLoading] = useState(false)
  const [wsStatus, setWsStatus] = useState<'disconnected' | 'connecting' | 'connected'>('disconnected')
  const [rbacInfo, setRbacInfo] = useState<any>(null)

  // Monitor inspectoría
  const [monitorCursos, setMonitorCursos] = useState<MonitorCurso[]>([])
  const [monitorAlumnosByCurso, setMonitorAlumnosByCurso] = useState<Record<string, Alumno[]>>({})
  const [monitorEventos, setMonitorEventos] = useState<Evento[]>([])
  const [monitorEstadosTemp, setMonitorEstadosTemp] = useState<EstadoTemporal[]>([])
  const [profLastByCurso, setProfLastByCurso] = useState<Record<string, number>>({})
  const [monitorSelectedCurso, setMonitorSelectedCurso] = useState<MonitorCurso | null>(null)
  const [monitorAlertas, setMonitorAlertas] = useState<Alerta[]>([])
  const [closingAlertaId, setClosingAlertaId] = useState<string | null>(null)

  useEffect(() => {
    const storedToken = localStorage.getItem('backoffice_token')
    if (storedToken) {
      setToken(storedToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
      cargarUsuario()
    } else {
      setLoading(false)
    }
  }, [])

  // WebSocket (realtime)
  useEffect(() => {
    if (!token) return

    const wsUrl = (() => {
      const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
      return `${proto}://${window.location.host}/ws`
    })()

    let ws: WebSocket | null = null
    let closed = false

    const connect = () => {
      if (closed) return
      setWsStatus('connecting')
      ws = new WebSocket(wsUrl)

      ws.onopen = () => setWsStatus('connected')
      ws.onclose = () => {
        setWsStatus('disconnected')
        if (!closed) {
          setTimeout(connect, 1500)
        }
      }
      ws.onerror = () => {
        // onclose will handle reconnect
      }
      ws.onmessage = (evt) => {
        const raw = String(evt.data || '')
        const parts = raw.split('\n').map(s => s.trim()).filter(Boolean)
        for (const p of parts) {
          try {
            const msg = JSON.parse(p)
            const type = msg?.type
            const payload = msg?.payload

            if (!type) continue

            // Eventos (creado/cerrado)
            if (type === 'evento_creado' && payload?.id) {
              setDashboard(prev => {
                if (!prev) return prev
                const ult = [payload, ...(prev.ultimos_eventos || [])].slice(0, 10)
                return { ...prev, eventos_activos: (prev.eventos_activos || 0) + 1, ultimos_eventos: ult }
              })
              setEventos(prev => {
                const exists = prev.some(e => e.id === payload.id)
                return exists ? prev : [payload, ...prev]
              })

              // monitor
              setMonitorEventos(prev => {
                const exists = prev.some(e => e.id === payload.id)
                return exists ? prev : [payload, ...prev]
              })
            }

            if (type === 'evento_cerrado' && payload?.id) {
              setDashboard(prev => {
                if (!prev) return prev
                return { ...prev, eventos_activos: Math.max(0, (prev.eventos_activos || 0) - 1) }
              })
              setEventos(prev => prev.filter(e => e.id !== payload.id))

              // monitor
              setMonitorEventos(prev => prev.filter(e => e.id !== payload.id))
            }

            // Estados temporales (solo para contador del dashboard)
            if (type === 'estado_temporal_creado') {
              setDashboard(prev => {
                if (!prev) return prev
                return { ...prev, estados_temporales_activos: (prev.estados_temporales_activos || 0) + 1 }
              })
              setMonitorEstadosTemp(prev => [payload, ...prev])
            }
            if (type === 'estado_temporal_cerrado') {
              setDashboard(prev => {
                if (!prev) return prev
                return { ...prev, estados_temporales_activos: Math.max(0, (prev.estados_temporales_activos || 0) - 1) }
              })
              const alumnoId = payload?.alumno_id
              if (alumnoId) {
                setMonitorEstadosTemp(prev => prev.map(e => e.alumno_id === alumnoId && !e.fin ? { ...e, fin: new Date().toISOString() } : e))
              }
            }

            // Acciones ejecutadas (para trazabilidad)
            if (type === 'accion_ejecutada' && payload?.exec?.id) {
              const exec = payload.exec
              const row: AccionEjecucion = {
                ...exec,
                accion: payload.accion,
                regla: payload.regla,
              }
              setAccionesEjecuciones(prev => [row, ...prev].slice(0, 50))
            }

            // Alertas operativas
            if (type === 'alerta_creada' && payload?.id) {
              setMonitorAlertas(prev => {
                const exists = prev.some(a => a.id === payload.id)
                return exists ? prev : [payload, ...prev].slice(0, 50)
              })
            }

            // Presencia profesor por bloque
            if (type === 'asistencia_bloque_registrada') {
              const cursoId = payload?.curso_id
              if (cursoId) {
                setProfLastByCurso(prev => ({ ...prev, [cursoId]: Date.now() }))
              }
            }
          } catch {
            // ignore parse errors
          }
        }
      }
    }

    connect()

    return () => {
      closed = true
      try {
        ws?.close()
      } catch {}
    }
  }, [token])

  const salaSemaforoForCursoId = (cursoId: string) => {
    const eventosCurso = monitorEventos.filter(e => e.curso_id === cursoId && e.activo)
    const estadosTempCurso = monitorEstadosTemp.filter(e => e.alumno?.curso_id === cursoId && !e.fin)
    const last = profLastByCurso[cursoId]

    // sin data
    if (eventosCurso.length === 0 && estadosTempCurso.length === 0 && !last) return 'gris'

    // rojo: SOS/COMPORTAMIENTO/DISCIPLINARIO
    const red = eventosCurso.some(e => ['SOS', 'COMPORTAMIENTO', 'DISCIPLINARIO'].includes(e?.concepto?.codigo || '')) ||
      estadosTempCurso.some(e => e.tipo === 'sos')
    if (red) return 'rojo'

    // amarillo: BANO/ENFERMERIA/INASISTENCIA
    const yellow = eventosCurso.some(e => ['BANO', 'ENFERMERIA', 'INASISTENCIA'].includes(e?.concepto?.codigo || '')) ||
      estadosTempCurso.some(e => e.tipo === 'bano' || e.tipo === 'enfermeria')
    if (yellow) return 'amarillo'

    return 'verde'
  }

  const salaColorForCursoId = (cursoId: string) => {
    const s = salaSemaforoForCursoId(cursoId)
    if (s === 'rojo') return '#dc2626'
    if (s === 'amarillo') return '#ca8a04'
    if (s === 'gris') return '#9ca3af'
    return '#16a34a'
  }

  const profColorForCursoId = (cursoId: string) => {
    const last = profLastByCurso[cursoId]
    if (!last) return '#9ca3af' // sin data
    const minutes = (Date.now() - last) / 60000
    if (minutes <= 75) return '#16a34a'
    if (minutes <= 150) return '#ca8a04'
    return '#9ca3af'
  }

  const studentColor = (alumnoId: string) => {
    // estados temporales activos
    const st = monitorEstadosTemp.find(e => e.alumno_id === alumnoId && !e.fin)
    if (st) {
      if (st.tipo === 'sos') return '#dc2626'
      return '#ca8a04'
    }
    // eventos activos (disciplina/comportamiento)
    const ev = monitorEventos.filter(e => e.alumno_id === alumnoId && e.activo)
    if (ev.some(e => ['SOS', 'COMPORTAMIENTO', 'DISCIPLINARIO'].includes(e?.concepto?.codigo || ''))) return '#dc2626'
    if (ev.some(e => ['BANO', 'ENFERMERIA'].includes(e?.concepto?.codigo || ''))) return '#ca8a04'
    return '#16a34a'
  }

  const cargarUsuario = async () => {
    try {
      const resp = await axios.get(`${API_URL}/auth/me`)
      setUsuario(resp.data)
      await cargarDatos()
    } catch (error) {
      console.error('Error cargando usuario:', error)
      logout()
    } finally {
      setLoading(false)
    }
  }

  const cargarDatos = async () => {
    await Promise.all([
      cargarDashboard(),
      cargarConceptos(),
      cargarAcciones(),
      cargarReglas(),
      cargarEventos()
    ])
  }

  const cargarDashboard = async () => {
    try {
      const resp = await axios.get(`${API_URL}/dashboard`)
      setDashboard(resp.data)
    } catch (error) {
      console.error('Error cargando dashboard:', error)
    }
  }

  const cargarConceptos = async () => {
    try {
      const resp = await axios.get(`${API_URL}/conceptos`)
      setConceptos(resp.data || [])
    } catch (error) {
      console.error('Error cargando conceptos:', error)
    }
  }

  const cargarAcciones = async () => {
    try {
      const resp = await axios.get(`${API_URL}/acciones`)
      setAcciones(resp.data || [])
    } catch (error) {
      console.error('Error cargando acciones:', error)
    }
  }

  const cargarReglas = async () => {
    try {
      const resp = await axios.get(`${API_URL}/reglas`)
      setReglas(resp.data || [])
    } catch (error) {
      console.error('Error cargando reglas:', error)
    }
  }

  const cargarEventos = async () => {
    try {
      const resp = await axios.get(`${API_URL}/eventos/activos`)
      setEventos(resp.data || [])
    } catch (error) {
      console.error('Error cargando eventos:', error)
    }
  }

  const cargarMonitor = async () => {
    try {
      const snapResp = await axios.get(`${API_URL}/monitor/snapshot`)
      const cursosData: MonitorCurso[] = snapResp.data?.cursos || []
      setMonitorCursos(cursosData)

      // inicializar profLastByCurso desde ultima_asistencia_en (para semáforo profesor consistente aunque llegues tarde)
      const base: Record<string, number> = {}
      for (const c of cursosData) {
        if (c.ultima_asistencia_en) {
          const t = Date.parse(c.ultima_asistencia_en)
          if (!Number.isNaN(t)) base[c.curso_id] = t
        }
      }
      setProfLastByCurso(prev => ({ ...base, ...prev }))

      // cargar eventos activos y estados temporales
      const [evtResp, estResp] = await Promise.all([
        axios.get(`${API_URL}/eventos/activos`),
        axios.get(`${API_URL}/estados-temporales`),
      ])
      setMonitorEventos(evtResp.data || [])
      setMonitorEstadosTemp(estResp.data || [])

      // alertas abiertas (cola)
      try {
        const a = await axios.get(`${API_URL}/alertas?estado=abierta&limit=50`)
        setMonitorAlertas(a.data || [])
      } catch {
        setMonitorAlertas([])
      }

      // precargar alumnos por curso (lazy: solo primeros N cursos para no explotar)
      const first = cursosData.slice(0, 20)
      const pairs = await Promise.all(
        first.map(async (c) => {
          const a = await axios.get(`${API_URL}/cursos/${c.curso_id}/alumnos`)
          return [c.curso_id, a.data || []] as const
        })
      )
      const map: Record<string, Alumno[]> = {}
      for (const [cid, arr] of pairs) map[cid] = arr
      setMonitorAlumnosByCurso(prev => ({ ...prev, ...map }))
    } catch (e) {
      console.error('Error cargando monitor:', e)
      setMonitorCursos([])
    }
  }

  const cerrarAlerta = async (id: string) => {
    setClosingAlertaId(id)
    try {
      await axios.put(`${API_URL}/alertas/${id}/cerrar`)
      setMonitorAlertas(prev => prev.filter(a => a.id !== id))
    } catch (e: any) {
      console.error(e)
      alert(e?.response?.data?.error || 'Error cerrando alerta')
    } finally {
      setClosingAlertaId(null)
    }
  }

  const cargarProfesores = async () => {
    try {
      const resp = await axios.get(`${API_URL}/usuarios?rol=profesor`)
      setProfesores(resp.data || [])
    } catch (error) {
      console.error('Error cargando profesores:', error)
      setProfesores([])
    }
  }

  const cargarCatalogosHorario = async () => {
    try {
      const [cursosResp, asignResp, bloquesResp] = await Promise.all([
        axios.get(`${API_URL}/cursos`),
        axios.get(`${API_URL}/asignaturas`),
        axios.get(`${API_URL}/bloques`),
      ])
      const cursosData = cursosResp.data || []
      setCursos(cursosData)
      setAsignaturas(asignResp.data || [])
      setBloques(bloquesResp.data || [])
      if (!cursoHorarioId && cursosData[0]?.id) {
        setCursoHorarioId(cursosData[0].id)
      }
    } catch (error) {
      console.error('Error cargando catálogos de horario:', error)
    }
  }

  const cargarHorariosCurso = async (cursoId: string) => {
    if (!cursoId) return
    try {
      const resp = await axios.get(`${API_URL}/horarios?curso_id=${cursoId}`)
      setHorarios(resp.data || [])
    } catch (error) {
      console.error('Error cargando horarios:', error)
      setHorarios([])
    }
  }

  const cargarTrazabilidad = async () => {
    setTraceLoading(true)
    try {
      const qs = auditTabla ? `?tabla=${encodeURIComponent(auditTabla)}&limit=50` : '?limit=50'
      const [audResp, execResp] = await Promise.all([
        axios.get(`${API_URL}/auditorias${qs}`),
        axios.get(`${API_URL}/acciones-ejecuciones?limit=50`),
      ])
      setAuditorias(audResp.data || [])
      setAccionesEjecuciones(execResp.data || [])
    } catch (error) {
      console.error('Error cargando trazabilidad:', error)
      setAuditorias([])
      setAccionesEjecuciones([])
    } finally {
      setTraceLoading(false)
    }
  }

  const cargarRBAC = async () => {
    try {
      const resp = await axios.get(`${API_URL}/auth/permisos`)
      setRbacInfo(resp.data)
    } catch (e) {
      console.error('Error cargando RBAC:', e)
      setRbacInfo(null)
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
      localStorage.setItem('backoffice_token', newToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      setToken(newToken)
      setUsuario(user)
      await cargarDatos()
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
      try {
        await axios.post(`${API_URL}/seed`)
      } catch (e) {}

      const resp = await axios.post(`${API_URL}/auth/login`, {
        email: 'backoffice@escuela.cl',
        password: 'profesor123'
      })
      const { token: newToken, usuario: user } = resp.data
      localStorage.setItem('backoffice_token', newToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      setToken(newToken)
      setUsuario(user)
      await cargarDatos()
    } catch (error: any) {
      console.error('Error en login demo:', error)
      alert(error?.response?.data?.error || 'Error en login demo')
    } finally {
      setLoginLoading(false)
    }
  }

  const logout = () => {
    localStorage.removeItem('backoffice_token')
    delete axios.defaults.headers.common['Authorization']
    setToken(null)
    setUsuario(null)
    setProfesores([])
    setCursos([])
    setAsignaturas([])
    setBloques([])
    setCursoHorarioId('')
    setHorarios([])
    setAuditorias([])
    setAccionesEjecuciones([])
  }

  const crearProfesor = async () => {
    if (!nuevoProfesorNombre || !nuevoProfesorEmail) {
      alert('Completa nombre y email')
      return
    }
    try {
      await axios.post(`${API_URL}/usuarios`, {
        nombre: nuevoProfesorNombre,
        email: nuevoProfesorEmail,
        rol: 'profesor',
        password: nuevoProfesorPassword || 'profesor123',
        activo: true,
      })
      setNuevoProfesorNombre('')
      setNuevoProfesorEmail('')
      setNuevoProfesorPassword('profesor123')
      await cargarProfesores()
    } catch (e: any) {
      console.error(e)
      alert(e?.response?.data?.error || 'Error creando profesor')
    }
  }

  const toggleProfesorActivo = async (u: Usuario) => {
    try {
      await axios.put(`${API_URL}/usuarios/${u.id}`, { activo: !(u.activo ?? true) })
      await cargarProfesores()
    } catch (e: any) {
      console.error(e)
      alert(e?.response?.data?.error || 'Error actualizando profesor')
    }
  }

  const resetProfesorPassword = async (u: Usuario) => {
    const newPass = prompt(`Nuevo password para ${u.email}:`, 'profesor123')
    if (!newPass) return
    try {
      await axios.put(`${API_URL}/usuarios/${u.id}`, { password: newPass })
      alert('Password actualizado')
    } catch (e: any) {
      console.error(e)
      alert(e?.response?.data?.error || 'Error actualizando password')
    }
  }

  const upsertHorario = async (dia: number, bloqueId: string, asignaturaId: string, profesorId: string) => {
    if (!cursoHorarioId) return
    const key = `${dia}:${bloqueId}`
    setSavingHorarioKey(key)
    try {
      await axios.post(`${API_URL}/horarios`, {
        curso_id: cursoHorarioId,
        dia_semana: dia,
        bloque_id: bloqueId,
        asignatura_id: asignaturaId,
        profesor_id: profesorId,
      })
      await cargarHorariosCurso(cursoHorarioId)
    } catch (e: any) {
      console.error(e)
      alert(e?.response?.data?.error || 'Error guardando horario')
    } finally {
      setSavingHorarioKey(null)
    }
  }

  const importarHorariosCSV = async () => {
    if (!csvImport.trim()) {
      alert('Pega un CSV primero')
      return
    }
    setCsvImportLoading(true)
    try {
      const resp = await axios.post(`${API_URL}/import/horarios`, { formato: 'csv', csv: csvImport })
      const d = resp.data
      alert(`Import OK: ${d.rows_ok} filas, errores: ${d.rows_error}, creados: ${d.horarios_creados}, actualizados: ${d.horarios_actualizados}`)
      if (cursoHorarioId) await cargarHorariosCurso(cursoHorarioId)
    } catch (e: any) {
      console.error(e)
      alert(e?.response?.data?.error || 'Error importando CSV')
    } finally {
      setCsvImportLoading(false)
    }
  }

  const cerrarEvento = async (eventoId: string) => {
    try {
      await axios.put(`${API_URL}/eventos/${eventoId}/cerrar`)
      await cargarEventos()
      await cargarDashboard()
    } catch (error: any) {
      console.error('Error cerrando evento:', error)
      alert(error?.response?.data?.error || 'Error cerrando evento')
    }
  }

  if (loading) {
    return (
      <div style={{ padding: '2rem', textAlign: 'center' }}>
        <p>Cargando...</p>
      </div>
    )
  }

  if (!token || !usuario) {
  return (
      <div style={{ padding: '2rem', textAlign: 'center', maxWidth: 400, margin: '0 auto' }}>
        <h1 style={{ marginBottom: '1rem' }}>Backoffice Escolar</h1>
        <p style={{ color: '#6b7280', marginBottom: '1.5rem' }}>
          Panel de administracion
        </p>

        <div style={{ display: 'grid', gap: '0.75rem' }}>
          <button
            onClick={loginDemo}
            disabled={loginLoading}
            style={{
              padding: '0.75rem 1.5rem',
              backgroundColor: '#7c3aed',
              color: 'white',
              border: 'none',
              borderRadius: 8,
              cursor: 'pointer',
              fontWeight: 700,
            }}
          >
            {loginLoading ? 'Ingresando...' : 'Login Demo (Backoffice)'}
          </button>

          <div style={{ color: '#9ca3af', fontSize: '0.875rem', margin: '0.5rem 0' }}>
            o ingresa manualmente:
          </div>

          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Email"
            style={{ padding: '0.75rem 1rem', borderRadius: 8, border: '1px solid #e5e7eb' }}
          />
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Contraseña"
            style={{ padding: '0.75rem 1rem', borderRadius: 8, border: '1px solid #e5e7eb' }}
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

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      {/* Sidebar */}
      <aside style={{
        width: 220,
        backgroundColor: '#1f2937',
        color: 'white',
        padding: '1rem',
      }}>
        <div style={{ marginBottom: '2rem' }}>
          <h2 style={{ fontSize: '1.1rem', fontWeight: 700 }}>Backoffice</h2>
          <p style={{ fontSize: '0.8rem', color: '#9ca3af' }}>{usuario.nombre}</p>
        </div>

        <nav style={{ display: 'grid', gap: '0.5rem' }}>
          {(['dashboard', 'monitor', 'profesores', 'horarios', 'trazabilidad', 'rbac', 'conceptos', 'acciones', 'reglas', 'eventos'] as Tab[]).map(t => (
            <button
              key={t}
              onClick={async () => {
                setTab(t)
                if (t === 'profesores') {
                  await cargarProfesores()
                }
                if (t === 'monitor') {
                  await cargarMonitor()
                }
                if (t === 'horarios') {
                  await Promise.all([cargarProfesores(), cargarCatalogosHorario()])
                  const id = cursoHorarioId || cursos[0]?.id
                  if (id) await cargarHorariosCurso(id)
                }
                if (t === 'trazabilidad') {
                  await cargarTrazabilidad()
                }
                if (t === 'rbac') {
                  await cargarRBAC()
                }
              }}
          style={{
                padding: '0.75rem 1rem',
                backgroundColor: tab === t ? '#374151' : 'transparent',
                color: 'white',
                border: 'none',
                borderRadius: 8,
                cursor: 'pointer',
                textAlign: 'left',
                fontWeight: tab === t ? 700 : 400,
              }}
            >
              {t === 'monitor'
                ? 'Monitor'
                : t === 'profesores'
                ? 'Profesores'
                : t === 'horarios'
                  ? 'Horarios'
                  : t === 'trazabilidad'
                    ? 'Trazabilidad'
                    : t === 'rbac'
                      ? 'RBAC'
                    : t.charAt(0).toUpperCase() + t.slice(1)}
            </button>
          ))}
        </nav>

        <button
          onClick={logout}
            style={{
            marginTop: '2rem',
            padding: '0.75rem 1rem',
            backgroundColor: '#ef4444',
            color: 'white',
            border: 'none',
            borderRadius: 8,
            cursor: 'pointer',
            width: '100%',
          }}
        >
          Cerrar sesion
        </button>
      </aside>

      {/* Main content */}
      <main style={{ flex: 1, padding: '2rem', backgroundColor: '#f9fafb' }}>
        {tab === 'rbac' && (
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
              <h1 style={{ margin: 0 }}>RBAC</h1>
              <button
                onClick={cargarRBAC}
                style={{ padding: '0.5rem 0.75rem', borderRadius: 8, border: '1px solid #e5e7eb', background: 'white', cursor: 'pointer' }}
              >
                Recargar
              </button>
            </div>

            <div style={{ background: 'white', border: '1px solid #e5e7eb', borderRadius: 12, padding: '1rem' }}>
              {!rbacInfo ? (
                <div style={{ color: '#6b7280' }}>Sin datos (carga RBAC).</div>
              ) : (
                <div style={{ display: 'grid', gap: '1rem' }}>
                  <div>
                    <div style={{ fontWeight: 900, marginBottom: 6 }}>Mi rol</div>
                    <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace' }}>{rbacInfo.rol}</div>
                  </div>
                  <div>
                    <div style={{ fontWeight: 900, marginBottom: 6 }}>Mis permisos</div>
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                      {(rbacInfo.mis_permisos || []).map((p: string) => (
                        <span key={p} style={{ fontSize: 12, background: '#f3f4f6', border: '1px solid #e5e7eb', padding: '4px 8px', borderRadius: 999 }}>
                          {p}
                        </span>
                      ))}
                    </div>
                  </div>
                  <div>
                    <div style={{ fontWeight: 900, marginBottom: 6 }}>Permisos por rol</div>
                    <pre style={{ margin: 0, fontSize: 12, padding: '0.75rem', borderRadius: 10, border: '1px solid #e5e7eb', background: '#111827', color: 'white', overflow: 'auto' }}>
{JSON.stringify(rbacInfo.por_rol || {}, null, 2)}
                    </pre>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {tab === 'monitor' && (
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
              <h1 style={{ margin: 0 }}>Monitor Inspectoría</h1>
              <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
                <div style={{ fontSize: 12, color: '#6b7280' }}>Realtime: <b>{wsStatus}</b></div>
                <button
                  onClick={cargarMonitor}
                  style={{ padding: '0.5rem 0.75rem', borderRadius: 8, border: '1px solid #e5e7eb', background: 'white', cursor: 'pointer' }}
                >
                  Recargar
                </button>
              </div>
            </div>

            <div style={{ background: 'white', border: '1px solid #e5e7eb', borderRadius: 12, padding: '0.75rem', marginBottom: '1rem' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                <div style={{ fontWeight: 900 }}>Alertas operativas (cola)</div>
                <div style={{ fontSize: 12, color: '#6b7280' }}>Abiertas: <b>{monitorAlertas.length}</b></div>
              </div>
              {monitorAlertas.length === 0 ? (
                <div style={{ color: '#6b7280' }}>Sin alertas abiertas.</div>
              ) : (
                <div style={{ display: 'grid', gap: 8 }}>
                  {monitorAlertas.slice(0, 10).map(a => (
                    <div key={a.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 12, padding: '0.6rem', border: '1px solid #f3f4f6', borderRadius: 10 }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                        <div
                          title={a.prioridad}
                          style={{
                            width: 10,
                            height: 10,
                            borderRadius: 999,
                            background:
                              a.prioridad === 'critica' ? '#b91c1c' :
                              a.prioridad === 'alta' ? '#dc2626' :
                              a.prioridad === 'media' ? '#ca8a04' : '#16a34a',
                          }}
                        />
                        <div>
                          <div style={{ fontWeight: 800 }}>{a.titulo}</div>
                          <div style={{ fontSize: 12, color: '#6b7280' }}>
                            {a.codigo} · {new Date(a.created_at).toLocaleString()}
                          </div>
                        </div>
                      </div>
                      <button
                        onClick={() => cerrarAlerta(a.id)}
                        disabled={closingAlertaId === a.id}
                        style={{ padding: '0.45rem 0.75rem', borderRadius: 8, border: 'none', background: '#111827', color: 'white', cursor: 'pointer', fontWeight: 800 }}
                      >
                        {closingAlertaId === a.id ? 'Cerrando...' : 'Atender/Cerrar'}
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '0.75rem' }}>
              {monitorCursos.map(c => (
                <button
                  key={c.curso_id}
                  onClick={async () => {
                    setMonitorSelectedCurso(c)
                    if (!monitorAlumnosByCurso[c.curso_id]) {
                      try {
                        const a = await axios.get(`${API_URL}/cursos/${c.curso_id}/alumnos`)
                        setMonitorAlumnosByCurso(prev => ({ ...prev, [c.curso_id]: a.data || [] }))
                      } catch {}
                    }
                  }}
                  style={{
                    textAlign: 'left',
                    background: 'white',
                    border: '1px solid #e5e7eb',
                    borderRadius: 12,
                    padding: '0.75rem',
                    cursor: 'pointer',
                  }}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '0.75rem' }}>
                    <div style={{ fontWeight: 800 }}>{c.nombre}</div>
                    <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                      <div title="Sala" style={{ width: 14, height: 14, borderRadius: 999, background: salaColorForCursoId(c.curso_id) }} />
                      <div title="Profesor" style={{ width: 14, height: 14, borderRadius: 999, background: profColorForCursoId(c.curso_id) }} />
                    </div>
                  </div>
                  <div style={{ fontSize: 12, color: '#6b7280', marginTop: 6 }}>
                    Eventos activos: {monitorEventos.filter(e => e.curso_id === c.curso_id && e.activo).length}
                  </div>
                </button>
              ))}
            </div>

            {monitorSelectedCurso && (
              <div
                style={{
                  position: 'fixed',
                  inset: 0,
                  background: 'rgba(0,0,0,0.35)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  padding: '1rem',
                }}
                onClick={() => setMonitorSelectedCurso(null)}
              >
                <div
                  style={{
                    width: 'min(900px, 96vw)',
                    maxHeight: '85vh',
                    overflow: 'auto',
                    background: 'white',
                    borderRadius: 14,
                    border: '1px solid #e5e7eb',
                    padding: '1rem',
                  }}
                  onClick={(e) => e.stopPropagation()}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.75rem' }}>
                    <div>
                      <div style={{ fontSize: 18, fontWeight: 900 }}>{monitorSelectedCurso.nombre}</div>
                      <div style={{ fontSize: 12, color: '#6b7280' }}>Semáforos: Sala / Profesor</div>
                    </div>
                    <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                      <div style={{ width: 14, height: 14, borderRadius: 999, background: salaColorForCursoId(monitorSelectedCurso.curso_id) }} />
                      <div style={{ width: 14, height: 14, borderRadius: 999, background: profColorForCursoId(monitorSelectedCurso.curso_id) }} />
                      <button onClick={() => setMonitorSelectedCurso(null)} style={{ marginLeft: 12, padding: '0.4rem 0.6rem', borderRadius: 8, border: '1px solid #e5e7eb', background: 'white', cursor: 'pointer' }}>
                        Cerrar
                      </button>
                    </div>
                  </div>

                  <div style={{ display: 'grid', gap: '0.5rem' }}>
                    {(monitorAlumnosByCurso[monitorSelectedCurso.curso_id] || []).map(a => (
                      <div key={a.id} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '0.75rem', padding: '0.6rem', border: '1px solid #f3f4f6', borderRadius: 10 }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                          <PersonIcon size={18} color={studentColor(a.id)} />
                          <div style={{ fontWeight: 700 }}>{a.apellido} {a.nombre}</div>
                        </div>
                        <div style={{ fontSize: 14 }}>
                          {monitorEventos.filter(e => e.alumno_id === a.id && e.activo).some(e => (e?.concepto?.codigo || '') === 'BANO') && '🚻 '}
                          {monitorEventos.filter(e => e.alumno_id === a.id && e.activo).some(e => (e?.concepto?.codigo || '') === 'ENFERMERIA') && '🏥 '}
                          {monitorEventos.filter(e => e.alumno_id === a.id && e.activo).some(e => (e?.concepto?.codigo || '') === 'SOS') && '🆘 '}
                          {monitorEventos.filter(e => e.alumno_id === a.id && e.activo).some(e => (e?.concepto?.codigo || '') === 'COMPORTAMIENTO') && '⚠️ '}
                        </div>
                      </div>
                    ))}
                    {(monitorAlumnosByCurso[monitorSelectedCurso.curso_id] || []).length === 0 && (
                      <div style={{ color: '#6b7280' }}>Sin alumnos cargados para este curso.</div>
                    )}
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {tab === 'trazabilidad' && (
              <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
              <h1 style={{ margin: 0 }}>Trazabilidad</h1>
              <button
                onClick={cargarTrazabilidad}
                disabled={traceLoading}
                style={{ padding: '0.5rem 0.75rem', borderRadius: 8, border: '1px solid #e5e7eb', background: 'white', cursor: 'pointer' }}
              >
                {traceLoading ? 'Cargando...' : 'Recargar'}
              </button>
                </div>

            <div style={{ background: 'white', padding: '1rem', borderRadius: 12, border: '1px solid #e5e7eb', marginBottom: '1rem' }}>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr auto', gap: '0.75rem', alignItems: 'end' }}>
                <div>
                  <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 6 }}>Filtrar auditoría por tabla (opcional)</div>
                  <input
                    value={auditTabla}
                    onChange={(e) => setAuditTabla(e.target.value)}
                    placeholder="ej: horarios, usuarios, eventos, asistencias..."
                    style={{ width: '100%', padding: '0.6rem', borderRadius: 8, border: '1px solid #e5e7eb' }}
                  />
              </div>
              <button
                  onClick={cargarTrazabilidad}
                  disabled={traceLoading}
                  style={{ padding: '0.7rem 1rem', borderRadius: 10, border: 'none', background: '#111827', color: 'white', cursor: 'pointer', fontWeight: 700 }}
                >
                  Aplicar
              </button>
              </div>
            </div>

            <div style={{ display: 'grid', gap: '1rem' }}>
              <div style={{ background: 'white', borderRadius: 12, border: '1px solid #e5e7eb', overflow: 'hidden' }}>
                <div style={{ padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', fontWeight: 700 }}>
                  Auditoría (últimos 50)
                    </div>
                <div style={{ padding: '0.75rem 1rem', overflowX: 'auto' }}>
                  {auditorias.length === 0 ? (
                    <p style={{ margin: 0, color: '#6b7280' }}>Sin registros.</p>
                  ) : (
                    <table style={{ width: '100%', borderCollapse: 'collapse', minWidth: 900 }}>
                      <thead>
                        <tr>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Fecha</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Tabla</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Acción</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Usuario</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Registro</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Detalle (preview)</th>
                        </tr>
                      </thead>
                      <tbody>
                        {auditorias.map(a => (
                          <tr key={a.id}>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6', whiteSpace: 'nowrap' }}>
                              {new Date(a.created_at).toLocaleString()}
                            </td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6' }}>{a.tabla}</td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6' }}>{a.accion}</td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6' }}>
                              {a.usuario?.email || a.usuario_id || '-'}
                            </td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace' }}>
                              {a.registro_id}
                            </td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace', fontSize: 12, color: '#374151' }}>
                              {JSON.stringify(a.datos_nuevos ?? a.datos_anteriores ?? {}).slice(0, 180)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )}
                    </div>
                  </div>

              <div style={{ background: 'white', borderRadius: 12, border: '1px solid #e5e7eb', overflow: 'hidden' }}>
                <div style={{ padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', fontWeight: 700 }}>
                  Acciones ejecutadas (reglas) - últimos 50
                </div>
                <div style={{ padding: '0.75rem 1rem', overflowX: 'auto' }}>
                  {accionesEjecuciones.length === 0 ? (
                    <p style={{ margin: 0, color: '#6b7280' }}>Sin ejecuciones registradas aún.</p>
                  ) : (
                    <table style={{ width: '100%', borderCollapse: 'collapse', minWidth: 900 }}>
                      <thead>
                        <tr>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Fecha</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Acción</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Regla</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Resultado</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Alumno</th>
                          <th style={{ textAlign: 'left', padding: '0.5rem', borderBottom: '1px solid #e5e7eb' }}>Detalle (preview)</th>
                        </tr>
                      </thead>
                      <tbody>
                        {accionesEjecuciones.map(e => (
                          <tr key={e.id}>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6', whiteSpace: 'nowrap' }}>
                              {new Date(e.ejecutado_en).toLocaleString()}
                            </td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6' }}>
                              {e.accion?.codigo || e.accion_id}
                            </td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6' }}>
                              {e.regla?.nombre || e.regla_id}
                            </td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6' }}>{e.resultado}</td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace', fontSize: 12 }}>
                              {e.alumno_id || '-'}
                            </td>
                            <td style={{ padding: '0.5rem', borderBottom: '1px solid #f3f4f6', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace', fontSize: 12, color: '#374151' }}>
                              {JSON.stringify(e.detalle ?? {}).slice(0, 180)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                )}
              </div>
              </div>
            </div>
          </div>
        )}

        {tab === 'profesores' && (
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
              <h1 style={{ margin: 0 }}>Profesores</h1>
              <button
                onClick={cargarProfesores}
                style={{ padding: '0.5rem 0.75rem', borderRadius: 8, border: '1px solid #e5e7eb', background: 'white', cursor: 'pointer' }}
              >
                Recargar
              </button>
                </div>

            <div style={{ background: 'white', padding: '1rem', borderRadius: 12, border: '1px solid #e5e7eb', marginBottom: '1rem' }}>
              <h3 style={{ marginTop: 0 }}>Crear profesor</h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr auto', gap: '0.75rem', alignItems: 'end' }}>
                <div>
                  <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 6 }}>Nombre</div>
                  <input value={nuevoProfesorNombre} onChange={(e) => setNuevoProfesorNombre(e.target.value)} style={{ width: '100%', padding: '0.6rem', borderRadius: 8, border: '1px solid #e5e7eb' }} />
                </div>
                <div>
                  <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 6 }}>Email</div>
                  <input value={nuevoProfesorEmail} onChange={(e) => setNuevoProfesorEmail(e.target.value)} style={{ width: '100%', padding: '0.6rem', borderRadius: 8, border: '1px solid #e5e7eb' }} />
                </div>
                <div>
                  <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 6 }}>Password inicial</div>
                  <input value={nuevoProfesorPassword} onChange={(e) => setNuevoProfesorPassword(e.target.value)} style={{ width: '100%', padding: '0.6rem', borderRadius: 8, border: '1px solid #e5e7eb' }} />
                </div>
                <button
                  onClick={crearProfesor}
                  style={{ padding: '0.7rem 1rem', borderRadius: 10, border: 'none', background: '#111827', color: 'white', cursor: 'pointer', fontWeight: 700 }}
                >
                  Crear
                </button>
                          </div>
              <div style={{ fontSize: 12, color: '#6b7280', marginTop: 10 }}>
                Password por defecto: <code>profesor123</code>
              </div>
            </div>

            <div style={{ background: 'white', borderRadius: 12, border: '1px solid #e5e7eb', overflow: 'hidden' }}>
              <div style={{ padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', fontWeight: 700 }}>Listado</div>
              <div style={{ padding: '0.75rem 1rem' }}>
                {profesores.length === 0 ? (
                  <p style={{ margin: 0, color: '#6b7280' }}>Sin profesores cargados (presiona “Recargar”).</p>
                ) : (
                  <div style={{ display: 'grid', gap: '0.5rem' }}>
                    {profesores.map(p => (
                      <div key={p.id} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr auto auto', gap: '0.75rem', alignItems: 'center', padding: '0.6rem', border: '1px solid #f3f4f6', borderRadius: 10 }}>
                        <div>
                          <div style={{ fontWeight: 700 }}>{p.nombre}</div>
                          <div style={{ fontSize: 12, color: '#6b7280' }}>{p.email}</div>
                        </div>
                        <div style={{ fontSize: 12, color: p.activo === false ? '#ef4444' : '#16a34a', fontWeight: 700 }}>
                          {p.activo === false ? 'INACTIVO' : 'ACTIVO'}
                        </div>
                        <button onClick={() => resetProfesorPassword(p)} style={{ padding: '0.45rem 0.75rem', borderRadius: 8, border: '1px solid #e5e7eb', background: 'white', cursor: 'pointer' }}>
                          Password
                        </button>
                        <button onClick={() => toggleProfesorActivo(p)} style={{ padding: '0.45rem 0.75rem', borderRadius: 8, border: 'none', background: '#f59e0b', color: 'white', cursor: 'pointer', fontWeight: 700 }}>
                          Activar/Desactivar
                        </button>
                      </div>
                            ))}
                          </div>
                )}
                        </div>
            </div>
                          </div>
                        )}

        {tab === 'horarios' && (
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
              <h1 style={{ margin: 0 }}>Horarios</h1>
              <div style={{ display: 'flex', gap: '0.5rem' }}>
                <button
                  onClick={async () => {
                    await Promise.all([cargarProfesores(), cargarCatalogosHorario()])
                    const id = cursoHorarioId || cursos[0]?.id
                    if (id) await cargarHorariosCurso(id)
                  }}
                  style={{ padding: '0.5rem 0.75rem', borderRadius: 8, border: '1px solid #e5e7eb', background: 'white', cursor: 'pointer' }}
                >
                  Recargar
                </button>
                      </div>
            </div>

            <div style={{ background: 'white', padding: '1rem', borderRadius: 12, border: '1px solid #e5e7eb', marginBottom: '1rem' }}>
              <h3 style={{ marginTop: 0 }}>Importar horarios (CSV)</h3>
              <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 8 }}>
                Formato: <code>curso,dia_semana,bloque_numero,asignatura,profesor_email</code> (día 1=Lun … 5=Vie)
              </div>
              <textarea
                value={csvImport}
                onChange={(e) => setCsvImport(e.target.value)}
                rows={6}
                style={{ width: '100%', padding: '0.75rem', borderRadius: 10, border: '1px solid #e5e7eb', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace' }}
              />
              <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 10 }}>
                <button
                  onClick={importarHorariosCSV}
                  disabled={csvImportLoading}
                  style={{ padding: '0.7rem 1rem', borderRadius: 10, border: 'none', background: '#111827', color: 'white', cursor: 'pointer', fontWeight: 700 }}
                >
                  {csvImportLoading ? 'Importando...' : 'Importar'}
                </button>
              </div>
            </div>

            <div style={{ background: 'white', padding: '1rem', borderRadius: 12, border: '1px solid #e5e7eb', marginBottom: '1rem' }}>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr auto', gap: '1rem', alignItems: 'end' }}>
                <div>
                  <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 6 }}>Curso</div>
                  <select
                    value={cursoHorarioId}
                    onChange={async (e) => {
                      const id = e.target.value
                      setCursoHorarioId(id)
                      await cargarHorariosCurso(id)
                    }}
                    style={{ width: '100%', padding: '0.6rem', borderRadius: 8, border: '1px solid #e5e7eb' }}
                  >
                    {(cursos || []).map(c => (
                      <option key={c.id} value={c.id}>{c.nombre}</option>
                    ))}
                  </select>
                </div>
                <button
                  onClick={() => cargarHorariosCurso(cursoHorarioId)}
                  style={{ padding: '0.7rem 1rem', borderRadius: 10, border: 'none', background: '#2563eb', color: 'white', cursor: 'pointer', fontWeight: 700 }}
                >
                  Ver
                </button>
              </div>
              <div style={{ fontSize: 12, color: '#6b7280', marginTop: 10 }}>
                Edita cada celda (asignatura y profesor). Se guarda al cambiar.
              </div>
            </div>

            {bloques.length === 0 || asignaturas.length === 0 || profesores.length === 0 ? (
              <div style={{ background: 'white', padding: '1rem', borderRadius: 12, border: '1px solid #e5e7eb' }}>
                <p style={{ margin: 0, color: '#6b7280' }}>
                  Faltan catálogos (bloques/asignaturas/profesores). Presiona “Recargar”.
                </p>
              </div>
            ) : (
              <div style={{ background: 'white', borderRadius: 12, border: '1px solid #e5e7eb', overflow: 'hidden' }}>
                <div style={{ padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', fontWeight: 700 }}>
                  Grilla (Lun–Vie) × (Bloques)
                </div>
                <div style={{ overflowX: 'auto' }}>
                  <table style={{ width: '100%', borderCollapse: 'collapse', minWidth: 900 }}>
                    <thead>
                      <tr>
                        <th style={{ textAlign: 'left', padding: '0.75rem', borderBottom: '1px solid #e5e7eb' }}>Bloque</th>
                        {[1, 2, 3, 4, 5].map(d => (
                          <th key={d} style={{ textAlign: 'left', padding: '0.75rem', borderBottom: '1px solid #e5e7eb' }}>
                            {d === 1 ? 'Lunes' : d === 2 ? 'Martes' : d === 3 ? 'Miércoles' : d === 4 ? 'Jueves' : 'Viernes'}
                          </th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {bloques.map(b => (
                        <tr key={b.id}>
                          <td style={{ padding: '0.75rem', borderBottom: '1px solid #f3f4f6', verticalAlign: 'top' }}>
                            <div style={{ fontWeight: 700 }}>#{b.numero}</div>
                            <div style={{ fontSize: 12, color: '#6b7280' }}>{b.hora_inicio}–{b.hora_fin}</div>
                          </td>
                          {[1, 2, 3, 4, 5].map(dia => {
                            const h = horarios.find(x => x.dia_semana === dia && x.bloque_id === b.id)
                            const key = `${dia}:${b.id}`
                            const saving = savingHorarioKey === key
                            const asignaturaId = h?.asignatura_id || asignaturas[0]?.id
                            const profesorId = h?.profesor_id || profesores[0]?.id
                            return (
                              <td key={dia} style={{ padding: '0.75rem', borderBottom: '1px solid #f3f4f6', verticalAlign: 'top' }}>
                                <div style={{ display: 'grid', gap: '0.5rem' }}>
                                  <select
                                    value={asignaturaId}
                                    onChange={(e) => upsertHorario(dia, b.id, e.target.value, profesorId)}
                                    disabled={saving}
                                    style={{ width: '100%', padding: '0.55rem', borderRadius: 8, border: '1px solid #e5e7eb' }}
                                  >
                                    {asignaturas.map(a => <option key={a.id} value={a.id}>{a.nombre}</option>)}
                                  </select>
                                  <select
                                    value={profesorId}
                                    onChange={(e) => upsertHorario(dia, b.id, asignaturaId, e.target.value)}
                                    disabled={saving}
                                    style={{ width: '100%', padding: '0.55rem', borderRadius: 8, border: '1px solid #e5e7eb' }}
                                  >
                                    {profesores.map(p => <option key={p.id} value={p.id}>{p.nombre}</option>)}
                                  </select>
                                  {saving && <div style={{ fontSize: 12, color: '#6b7280' }}>Guardando...</div>}
                                </div>
                              </td>
                            )
                          })}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
            </div>
        )}

        {tab === 'dashboard' && dashboard && (
          <div>
            <h1 style={{ marginBottom: '1.5rem' }}>Dashboard</h1>
            <div style={{ marginBottom: '1rem', fontSize: 12, color: '#6b7280' }}>
              Realtime: <span style={{ fontWeight: 700 }}>{wsStatus}</span>
            </div>

            <div style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
              gap: '1rem',
              marginBottom: '2rem',
            }}>
              <div style={{ background: 'white', padding: '1.5rem', borderRadius: 12, border: '1px solid #e5e7eb' }}>
                <div style={{ fontSize: '2rem', fontWeight: 700 }}>{dashboard.total_cursos}</div>
                <div style={{ color: '#6b7280' }}>Cursos</div>
              </div>
              <div style={{ background: 'white', padding: '1.5rem', borderRadius: 12, border: '1px solid #e5e7eb' }}>
                <div style={{ fontSize: '2rem', fontWeight: 700 }}>{dashboard.total_alumnos}</div>
                <div style={{ color: '#6b7280' }}>Alumnos</div>
              </div>
              <div style={{ background: 'white', padding: '1.5rem', borderRadius: 12, border: '1px solid #e5e7eb' }}>
                <div style={{ fontSize: '2rem', fontWeight: 700, color: '#ef4444' }}>{dashboard.eventos_activos}</div>
                <div style={{ color: '#6b7280' }}>Eventos Activos</div>
              </div>
              <div style={{ background: 'white', padding: '1.5rem', borderRadius: 12, border: '1px solid #e5e7eb' }}>
                <div style={{ fontSize: '2rem', fontWeight: 700, color: '#f59e0b' }}>{dashboard.estados_temporales_activos}</div>
                <div style={{ color: '#6b7280' }}>Estados Temporales</div>
              </div>
            </div>

            <div style={{ background: 'white', padding: '1.5rem', borderRadius: 12, border: '1px solid #e5e7eb', marginBottom: '2rem' }}>
              <h3 style={{ marginBottom: '1rem' }}>Asistencia Hoy</h3>
              <div style={{ display: 'flex', gap: '2rem' }}>
                <div>
                  <span style={{ fontSize: '1.5rem', fontWeight: 700, color: '#22c55e' }}>{dashboard.asistencia_hoy.presentes}</span>
                  <span style={{ color: '#6b7280', marginLeft: '0.5rem' }}>Presentes</span>
                </div>
                <div>
                  <span style={{ fontSize: '1.5rem', fontWeight: 700, color: '#ef4444' }}>{dashboard.asistencia_hoy.ausentes}</span>
                  <span style={{ color: '#6b7280', marginLeft: '0.5rem' }}>Ausentes</span>
                </div>
                <div>
                  <span style={{ fontSize: '1.5rem', fontWeight: 700, color: '#3b82f6' }}>{dashboard.asistencia_hoy.justificados}</span>
                  <span style={{ color: '#6b7280', marginLeft: '0.5rem' }}>Justificados</span>
                </div>
              </div>
          </div>
        </div>
      )}

        {tab === 'conceptos' && (
          <div>
            <h1 style={{ marginBottom: '1.5rem' }}>Conceptos</h1>
            <p style={{ color: '#6b7280', marginBottom: '1rem' }}>
              Los conceptos son eventos estandarizados que pueden ocurrir en el establecimiento.
            </p>

            <div style={{ display: 'grid', gap: '0.75rem' }}>
              {conceptos.map(c => (
                <div key={c.id} style={{
                  background: 'white',
                  padding: '1rem',
                  borderRadius: 8,
                  border: '1px solid #e5e7eb',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}>
                  <div>
                    <div style={{ fontWeight: 700 }}>{c.nombre}</div>
                    <div style={{ fontSize: '0.85rem', color: '#6b7280' }}>Codigo: {c.codigo}</div>
                    {c.descripcion && <div style={{ fontSize: '0.85rem', color: '#9ca3af' }}>{c.descripcion}</div>}
          </div>
                  <div style={{
                    padding: '0.25rem 0.75rem',
                    borderRadius: 999,
                    backgroundColor: c.activo ? '#dcfce7' : '#fef2f2',
                    color: c.activo ? '#166534' : '#991b1b',
                    fontSize: '0.8rem',
                    fontWeight: 600,
                  }}>
                    {c.activo ? 'Activo' : 'Inactivo'}
            </div>
          </div>
              ))}
              {conceptos.length === 0 && (
                <p style={{ color: '#6b7280' }}>No hay conceptos configurados</p>
              )}
            </div>
          </div>
        )}

        {tab === 'acciones' && (
          <div>
            <h1 style={{ marginBottom: '1.5rem' }}>Acciones</h1>
            <p style={{ color: '#6b7280', marginBottom: '1rem' }}>
              Las acciones son respuestas automaticas que se ejecutan cuando se disparan reglas.
            </p>

            <div style={{ display: 'grid', gap: '0.75rem' }}>
              {acciones.map(a => (
                <div key={a.id} style={{
                  background: 'white',
                  padding: '1rem',
                  borderRadius: 8,
                  border: '1px solid #e5e7eb',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}>
                  <div>
                    <div style={{ fontWeight: 700 }}>{a.nombre}</div>
                    <div style={{ fontSize: '0.85rem', color: '#6b7280' }}>
                      Codigo: {a.codigo} | Tipo: {a.tipo}
            </div>
          </div>
                  <div style={{
                    padding: '0.25rem 0.75rem',
                    borderRadius: 999,
                    backgroundColor: a.activo ? '#dcfce7' : '#fef2f2',
                    color: a.activo ? '#166534' : '#991b1b',
                    fontSize: '0.8rem',
                    fontWeight: 600,
                  }}>
                    {a.activo ? 'Activo' : 'Inactivo'}
        </div>
                </div>
              ))}
              {acciones.length === 0 && (
                <p style={{ color: '#6b7280' }}>No hay acciones configuradas</p>
              )}
            </div>
          </div>
        )}

        {tab === 'reglas' && (
          <div>
            <h1 style={{ marginBottom: '1.5rem' }}>Reglas</h1>
            <p style={{ color: '#6b7280', marginBottom: '1rem' }}>
              Las reglas definen condiciones que, al cumplirse, disparan acciones automaticas.
            </p>

            <div style={{ display: 'grid', gap: '0.75rem' }}>
              {reglas.map(r => (
                <div key={r.id} style={{
                  background: 'white',
                  padding: '1rem',
                  borderRadius: 8,
                  border: '1px solid #e5e7eb',
                }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                    <div>
                      <div style={{ fontWeight: 700 }}>{r.nombre}</div>
                      <div style={{ fontSize: '0.85rem', color: '#6b7280', marginTop: '0.25rem' }}>
                        Concepto: {r.concepto?.nombre || r.concepto_id}
                    </div>
                      <div style={{ fontSize: '0.85rem', color: '#6b7280' }}>
                        Accion: {r.accion?.nombre || r.accion_id}
                      </div>
                      <div style={{ fontSize: '0.8rem', color: '#9ca3af', marginTop: '0.5rem' }}>
                        Condicion: {JSON.stringify(r.condicion)}
                      </div>
                    </div>
                    <div style={{
                      padding: '0.25rem 0.75rem',
                      borderRadius: 999,
                      backgroundColor: r.activo ? '#dcfce7' : '#fef2f2',
                      color: r.activo ? '#166534' : '#991b1b',
                      fontSize: '0.8rem',
                      fontWeight: 600,
                    }}>
                      {r.activo ? 'Activo' : 'Inactivo'}
                    </div>
                  </div>
                </div>
              ))}
              {reglas.length === 0 && (
                <p style={{ color: '#6b7280' }}>No hay reglas configuradas</p>
              )}
            </div>
          </div>
        )}

        {tab === 'eventos' && (
          <div>
            <h1 style={{ marginBottom: '1.5rem' }}>Eventos Activos</h1>

            <div style={{ display: 'grid', gap: '0.75rem' }}>
              {eventos.map(e => (
                <div key={e.id} style={{
                  background: 'white',
                  padding: '1rem',
                  borderRadius: 8,
                  border: '1px solid #e5e7eb',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}>
                  <div>
                    <div style={{ fontWeight: 700 }}>{e.concepto?.nombre || 'Evento'}</div>
                    <div style={{ fontSize: '0.85rem', color: '#6b7280' }}>
                      {e.alumno && `Alumno: ${e.alumno.nombre} ${e.alumno.apellido}`}
                      {e.curso && ` | Curso: ${e.curso.nombre}`}
                    </div>
                    <div style={{ fontSize: '0.8rem', color: '#9ca3af' }}>
                      Origen: {e.origen} | {new Date(e.created_at).toLocaleString('es-CL')}
                    </div>
                  </div>
                  <button
                    onClick={() => cerrarEvento(e.id)}
                    style={{
                      padding: '0.5rem 1rem',
                      backgroundColor: '#22c55e',
                      color: 'white',
                      border: 'none',
                      borderRadius: 8,
                      cursor: 'pointer',
                      fontWeight: 600,
                    }}
                  >
                    Cerrar
                  </button>
              </div>
            ))}
              {eventos.length === 0 && (
                <p style={{ color: '#6b7280' }}>No hay eventos activos</p>
              )}
              </div>
                </div>
              )}
    </main>
    </div>
  )
}
