# Modelo de Conceptos, Acciones y Reglas

## Conceptos
Un Concepto representa un evento estandarizado del establecimiento.
Los conceptos no ejecutan logica por si mismos; actuan como disparadores.

**Ejemplos:**
- Alumno necesita ir al bano
- Profesor inicia la clase
- Alumno ausente
- Uso indebido de celular
- Incidente disciplinario

## Acciones
Respuestas automaticas configurables asociadas a uno o mas conceptos.

**Ejemplos:**
- Cambiar estado visual en dashboard
- Enviar notificacion a apoderado
- Crear solicitud para inspectoria
- Registrar anotacion en libro digital
- Generar alerta a Asistente Social

Las acciones pueden activarse, desactivarse y parametrizarse desde backoffice.

## Reglas
Condiciones temporales o contextuales que disparan acciones.

**Ejemplos:**
- 2 inasistencias sin justificativo medico en 7 dias
- 3 inasistencias no consecutivas en 30 dias
- Alumno marcado como caso especial inhibe alertas automaticas

Las reglas se evaluan mediante un motor de reglas independiente.

## Eventos
Ocurrencia concreta de un concepto para un alumno, curso o clase.

Cada evento:
- Tiene timestamp
- Tiene origen (profesor, sistema)
- Ejecuta acciones
- Queda registrado en auditoria
