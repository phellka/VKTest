import React, { useEffect, useState } from 'react';
import axios from 'axios';
import 'bootstrap/dist/css/bootstrap.min.css';

const ContainersTable = () => {
  
  const [containers, setContainers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const formatDate = (isoDate) => {
    if (!isoDate) { 
      return "";
    }
    const date = new Date(isoDate);
    const day = String(date.getUTCDate()).padStart(2, '0');
    const month = String(date.getUTCMonth() + 1).padStart(2, '0'); 
    const year = date.getFullYear();
    const hours = String(date.getUTCHours()).padStart(2, '0');
    const minutes = String(date.getUTCMinutes()).padStart(2, '0');
    const seconds = String(date.getUTCSeconds()).padStart(2, '0');

    return `${day}.${month}.${year} - ${hours}:${minutes}:${seconds}`;
  };

  const getCurrentDateTime = () => {
    const currentDate = new Date();
    return formatDate(currentDate);
  };

  const fetchContainers = () => {
    setLoading(true);
    axios
      .get("/api/containers/with-last-ping")
      .then((response) => {
        const updatedContainers = response.data?.map((container) => ({
          ...container,
          lastSuccessfulPingDate: formatDate(container.Timestamp),
        })) || [];
        
        setContainers(updatedContainers);
        setLoading(false);
      })
      .catch((error) => {
        setError(error);
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchContainers();

    const intervalId = setInterval(() => {
      fetchContainers();
    }, 15000);

    return () => clearInterval(intervalId); 
  }, []);



  if (loading) return <div>Загрузка...</div>;

  if (error) return <div>{error}</div>;

  if (containers.length === 0) {
    return (
      <div className="container mt-4">
        <h1>Список контейнеров</h1>
        <div className="mb-4">
          <strong>Текущая дата и время UTC:</strong> {getCurrentDateTime()}
        </div>
        <div>Контейнеры не найдены.</div>
      </div>
    );
  }

  return (
    <div className="container mt-4">
      <h1>Список контейнеров</h1>
      <div className="mb-4">
        <strong>Текущая дата и время UTC:</strong> {getCurrentDateTime()}
      </div>
      <table className="table table-bordered table-striped mt-4">
        <thead className="thead-dark">
          <tr>
            <th>Имя</th>
            <th>IP</th>
            <th>Длительность пинга, мс.</th>
            <th>Последняя успешная попытка UTC</th>
          </tr>
        </thead>
        <tbody>
          {containers.map((container) => {
            const rowClass = container.lastSuccessfulPingDate ? 'table-success' : 'table-danger';
            return (
              <tr key={container.ID} className={rowClass}>
                <td>{container.Name}</td>
                <td>{container.Ip}</td>
                <td>{container.Pingtime}</td>
                <td>{container.lastSuccessfulPingDate ? container.lastSuccessfulPingDate : ''}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

export default ContainersTable;
