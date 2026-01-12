import React, { useEffect, useState } from 'react';
import {
  SafeAreaView,
  ScrollView,
  StyleSheet,
  Text,
  View,
  TouchableOpacity,
  Alert,
  RefreshControl,
} from 'react-native';
import axios from 'axios';const API_URL = 'http://localhost:8080/api/v1';interface Evento {
  id: string;
  sala_id: string;
  sala_numero?: string;
  tipo: string;
  descripcion: string;
  estado: 'verde' | 'amarillo' | 'rojo';
  prioridad: 'baja' | 'media' | 'alta' | 'critica';
  hora_inicio: string;
}function App(): JSX.Element {
  const [eventos, setEventos] = useState<Evento[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);  useEffect(() => {
    cargarEventos();
  }, []);  const cargarEventos = async () => {
    try {
      const response = await axios.get(`${API_URL}/eventos`);
      setEventos(response.data);
    } catch (error) {
      console.error('Error cargando eventos:', error);
      Alert.alert('Error', 'No se pudieron cargar los eventos');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };  const cerrarEvento = async (eventoId: string) => {
    Alert.alert(
      'Cerrar Evento',
      '¿Confirmar que el evento ha sido resuelto?',
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Confirmar',
          onPress: async () => {
            try {
              await axios.put(`${API_URL}/eventos/${eventoId}/cerrar`, {
                accion_tomada: 'Evento resuelto en terreno',
              });
              await cargarEventos();
              Alert.alert('Éxito', 'Evento cerrado correctamente');
            } catch (error) {
              console.error('Error cerrando evento:', error);
              Alert.alert('Error', 'No se pudo cerrar el evento');
            }
          },
        },
      ]
    );
  };  const getSemaforoEmoji = (estado: string) => {
    switch (estado) {
      case 'verde':
        return '🟢';
      case 'amarillo':
        return '🟡';
      case 'rojo':
        return '🔴';
      default:
        return '⚪';
    }
  };  const getPrioridadColor = (prioridad: string): string => {
    switch (prioridad) {
      case 'critica':
        return '#ef4444';
      case 'alta':
        return '#f59e0b';
      case 'media':
        return '#eab308';
      case 'baja':
        return '#3b82f6';
      default:
        return '#6b7280';
    }
  };  const onRefresh = () => {
    setRefreshing(true);
    cargarEventos();
  };  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Alertas Inspectoría</Text>
        <Text style={styles.subtitle}>
          {eventos.length} evento(s) activo(s)
        </Text>
      </View>      <ScrollView
        style={styles.scrollView}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }>
        {loading ? (
          <View style={styles.centerContainer}>
            <Text>Cargando eventos...</Text>
          </View>
        ) : eventos.length === 0 ? (
          <View style={styles.centerContainer}>
            <Text style={styles.emptyText}>No hay eventos activos</Text>
          </View>
        ) : (
          eventos.map((evento) => (
            <View
              key={evento.id}
              style={[
                styles.eventoCard,
                {
                  borderLeftColor: getPrioridadColor(evento.prioridad),
                  borderLeftWidth: 4,
                },
              ]}>
              <View style={styles.eventoHeader}>
                <View style={styles.eventoTitle}>
                  <Text style={styles.semaforo}>
                    {getSemaforoEmoji(evento.estado)}
                  </Text>
                  <Text style={styles.salaText}>
                    Sala {evento.sala_numero || evento.sala_id}
                  </Text>
                </View>
                <View
                  style={[
                    styles.prioridadBadge,
                    { backgroundColor: getPrioridadColor(evento.prioridad) },
                  ]}>
                  <Text style={styles.prioridadText}>
                    {evento.prioridad.toUpperCase()}
                  </Text>
                </View>
              </View>              <Text style={styles.tipoText}>Tipo: {evento.tipo}</Text>
              <Text style={styles.descripcionText}>{evento.descripcion}</Text>              <Text style={styles.horaText}>
                Inicio: {new Date(evento.hora_inicio).toLocaleString('es-CL')}
              </Text>              <TouchableOpacity
                style={styles.cerrarButton}
                onPress={() => cerrarEvento(evento.id)}>
                <Text style={styles.cerrarButtonText}>Marcar como Resuelto</Text>
              </TouchableOpacity>
            </View>
          ))
        )}
      </ScrollView>
    </SafeAreaView>
  );
}const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    backgroundColor: '#ffffff',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
  },
  subtitle: {
    fontSize: 14,
    color: '#6b7280',
    marginTop: 4,
  },
  scrollView: {
    flex: 1,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  emptyText: {
    fontSize: 16,
    color: '#6b7280',
  },
  eventoCard: {
    backgroundColor: '#ffffff',
    margin: 12,
    padding: 16,
    borderRadius: 8,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
  eventoHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  eventoTitle: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  semaforo: {
    fontSize: 20,
    marginRight: 8,
  },
  salaText: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#111827',
  },
  prioridadBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
  },
  prioridadText: {
    color: '#ffffff',
    fontSize: 10,
    fontWeight: 'bold',
  },
  tipoText: {
    fontSize: 14,
    color: '#6b7280',
    marginBottom: 4,
  },
  descripcionText: {
    fontSize: 16,
    color: '#111827',
    marginBottom: 8,
  },
  horaText: {
    fontSize: 12,
    color: '#9ca3af',
    marginBottom: 12,
  },
  cerrarButton: {
    backgroundColor: '#22c55e',
    padding: 12,
    borderRadius: 6,
    alignItems: 'center',
  },
  cerrarButtonText: {
    color: '#ffffff',
    fontSize: 16,
    fontWeight: '600',
  },
});export default App;
