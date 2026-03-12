# Guía de Observabilidad y Monitoreo - Prunus API

Esta guía detalla cómo integrar Prunus con sistemas de agregación de logs y monitoreo como **Grafana Loki**, **AWS CloudWatch** y **Fluentd**, aprovechando el estándar de *Structured Logging* (JSON) implementado en la aplicación.

## 1. Estrategia de Logging Estructurado

Prunus utiliza `log/slog` de Go para generar logs en formato JSON. Este formato es el estándar de la industria para la observabilidad moderna, ya que permite que herramientas externas parseen y analicen los datos sin necesidad de expresiones regulares complejas.

### Ejemplo de Log de Rendimiento
```json
{
  "time": "2026-03-11T14:30:00Z",
  "level": "WARN",
  "msg": "Operación lenta detectada",
  "capa": "store",
  "operacion": "GetAllUsuarios",
  "latencia_ms": 250,
  "threshold": "200ms"
}
```

---

## 2. Integración con Grafana Loki

Loki es el sistema de agregación de logs preferido para el ecosistema Grafana. Se recomienda el uso de **Promtail** para recolectar los logs de Prunus.

### Configuración de Promtail (recolector)
Promtail debe configurarse para leer los logs de la salida estándar (si se usa Docker) o del archivo de log definido en el servicio systemd (`/var/log/prunus/output.log`).

```yaml
scrape_configs:
  - job_name: prunus-api
    static_configs:
      - targets:
          - localhost
        labels:
          job: prunus
          env: production
          __path__: /var/log/prunus/*.log
    pipeline_stages:
      - json:
          expressions:
            level: level
            msg: msg
            latencia: latencia_ms
```

---

## 3. Integración con AWS CloudWatch

Si el despliegue se realiza en AWS (EC2 o ECS), CloudWatch Logs capturará automáticamente el formato JSON.

### Beneficios en CloudWatch:
- **CloudWatch Logs Insights**: Permite ejecutar consultas SQL-like sobre los logs.
- **Métricas de Filtro**: Puede crear alarmas basadas en el campo `latencia_ms`.

**Ejemplo de Consulta Insights:**
```sql
fields @timestamp, msg, latencia_ms
| filter latencia_ms > 200
| sort @timestamp desc
```

---

## 4. Integración con Fluentd (Log Forwarder)

Fluentd actúa como una capa unificadora para enviar logs a múltiples destinos (Elasticsearch, S3, Datadog).

### Configuración de Entrada (Input)
```conf
<source>
  @type tail
  path /var/log/prunus/output.log
  pos_file /var/log/prunus/output.log.pos
  tag prunus.api
  <parse>
    @type json
    time_key time
    time_format %Y-%m-%dT%H:%M:%S.%LZ
  </parse>
</source>
```

---

## 5. Visualización en Grafana

Con los logs indexados en Loki o CloudWatch, puede crear tableros (Dashboards) interactivos:

1. **Gráfico de Latencia**: Visualice la evolución de `latencia_ms` por operación.
2. **Mapa de Calor de Errores**: Identifique picos de logs con `level: "ERROR"`.
3. **Métricas de Negocio**: Si el log incluye el ID de la sucursal, puede monitorear la actividad por ubicación en tiempo real.

---

## Recomendaciones para Producción

1. **Retención**: Configure políticas de retención de logs (ej. 7-14 días) para evitar costos excesivos de almacenamiento.
2. **Nivel de Log**: En producción, utilice `level: INFO` para reducir el ruido, activando `DEBUG` solo durante investigaciones de fallos.
3. **Alertas**: Configure alertas automáticas si el número de errores (`level: ERROR`) supera un umbral en un periodo de 5 minutos.
