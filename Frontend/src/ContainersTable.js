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

  useEffect(() => {
    axios
      .get(`${process.env.REACT_APP_BACKEND_URL}/containers/with-last-ping`)
      .then(async (response) => {
        const updatedContainers = response.data.map((container) => ({
          ...container,
          lastSuccessfulPingDate:  formatDate(container.Timestamp),
        }));
        setContainers(updatedContainers);
        setLoading(false);
      })
      .catch((error) => {
        setError('Ошибка при загрузке данных');
        console.log(process.env.REACT_APP_BACKEND_URL);
        setLoading(false);
      });
  }, []);


  if (loading) return <div>Загрузка...</div>;


  if (error) return <div>{error}</div>;

  return (
    <div className="container mt-4">
      <h1>Список контейнеров</h1>
      <div className="mb-4">
        <strong>Текущая дата и время UTC:</strong> {getCurrentDateTime()}
      </div>
      <table className="table table-bordered table-striped mt-4">
        <thead className="thead-dark">
          <tr>
            <th>ID</th>
            <th>Имя</th>
            <th>IP</th>
            <th>Последняя успешная попытка UTC</th>
          </tr>
        </thead>
        <tbody>
        {containers.map((container) => {
            const rowClass = container.lastSuccessfulPingDate ? 'table-success' : 'table-danger';
            return (
              <tr key={container.ID} className={rowClass}>
                <td>{container.ID}</td>
                <td>{container.Name}</td>
                <td>{container.Ip}</td>
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
