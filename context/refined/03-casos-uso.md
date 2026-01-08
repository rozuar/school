# Casos de Uso

## Registro de Asistencia

### Flujo Principal
1. El profesor inicia la clase (bloque horario)
2. Registra alumnos presentes y ausentes
3. El sistema genera eventos de inasistencia
4. Se notifica a apoderados
5. Se habilita canal de justificacion

### Asistencia por Bloque/Ramo
- La asistencia se toma por cada bloque horario
- Cuando un curso cambia de ramo/profesor, se debe retomar la asistencia
- Cada profesor es responsable de la asistencia de su bloque
- Un alumno puede estar presente en un bloque y ausente en otro

### Estados del Alumno Durante Clase
- Presente
- Ausente
- En bano/enfermeria (estado temporal)
- SOS/Comportamiento grave

## Justificacion de Inasistencia
1. Apoderado recibe notificacion
2. Puede adjuntar certificado medico o justificacion libre
3. El estado de la inasistencia cambia
4. Se actualiza el registro del alumno

## Seguimiento Social
1. Reglas detectan inasistencias reiteradas
2. Se genera alerta a Asistente Social
3. El caso queda en monitoreo
4. Se registra seguimiento y acciones tomadas

## Casos Especiales
- Alumnos pueden marcarse como casos especiales
- Inhiben alertas automaticas
- Siguen siendo monitoreados
- Requieren seguimiento personalizado

## Gestion de Permisos Temporales
- Alumno solicita ir al bano
- Profesor activa concepto "Bano"
- Estado visual cambia a amarillo
- Al regresar, profesor desactiva el concepto
- Tiempo de ausencia queda registrado
