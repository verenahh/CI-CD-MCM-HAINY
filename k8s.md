# Task 4: Production Readiness
## 2. Health Checks

### Readiness Probe vs Liveness Probe & what happens when it fails
Eine "readiness probe" legt fest ob ein container traffic empfangen kann oder nicht. Wenn dieser Test fehlschlägt, läuft der pod zwar weiter aber wird von service endpoints entfernt damit er kein weiterer traffic hingeschickt wird.

Eine "liveness probe" stellt fest ob ein container weiterhin korrekt läuft. Wenn dieser Test fehlschlägt wird der container von Kubernetes neu gestartet.

### Why different initialDelaySeconds values are used

Unterschiedliche "initialDelaySeconds" werden verwendet weil unterschiedliche Appliaktionen unterschiedliche Zeit brauchen um gestartet zu werden und stabil zu laufen. Eine readiness probe ist kürzer weil sie nur kurz checkt ob man traffic auf die applikation routen kann, die liveness probe soll im Idealfall länger sein um zu vermeiden dass der Test failed nur weil der container noch gar nicht fertig hochgefahren ist.

---

## 3. Resource Limits

### What happens if CPU or memory limits are exceeded

Wenn ein Container sein CPU-Limit überschreitet, wird seine Auslastung gedrosselt, was die Ausführung verlangsamt aber den Container nicht gleich stoppt. Wenn ein Container sein Speicherlimit (memory limit) überschreitet, wird er vom System beendet und kann je nach Neustartrichtlinie von Kubernetes neu gestartet werden.

### Why specify both requests and limits

Ressourcenanforderungen legen die Mindestmenge an CPU-Leistung und Arbeitsspeicher fest, die einem Container garantiert wird. Kubernetes nutzt diese Informationen, um zu entscheiden, wo der Pod eingeplant wird.

Ressourcenbeschränkungen legen die maximale Menge an CPU-Leistung und Arbeitsspeicher fest, die ein Container nutzen darf. Dadurch wird verhindert, dass ein einzelner Container übermäßig viele Ressourcen beansprucht und andere Workloads im Cluster beeinträchtigt.

Die Verwendung beider Mechanismen gewährleistet eine stabile Einplanung und eine kontrollierte Ressourcennutzung innerhalb des Clusters.