import React, { useEffect, useState } from 'react';
import axios from 'axios';
import 'bootstrap/dist/css/bootstrap.min.css';

const ContainersTable = () => {
  
  const [containers, setContainers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchLastSuccessfulPingDate = async (containerId) => {
    try {
      const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/container/lastsuccessful?id=${containerId}`);
      return formatDate(response.data.Timestamp);
    } catch (err) {
      if (err.response && err.response.status === 404) {
        return null;
      }
      return null; 
    }
  };

  const formatDate = (isoDate) => {
    const date = new Date(isoDate);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0'); 
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    return `${day}.${month}.${year} - ${hours}:${minutes}:${seconds}`;
  };

  const getCurrentDateTime = () => {
    const currentDate = new Date();
    return formatDate(currentDate.toISOString());
  };

  useEffect(() => {
    axios
      .get(`${process.env.REACT_APP_BACKEND_URL}/containers`)
      .then(async (response) => {
        const updatedContainers = await Promise.all(
          response.data.map(async (container) => {
            const lastSuccessfulPingDate = container.LastSuccessfulPingId
              ? await fetchLastSuccessfulPingDate(container.ID)
              : null; 
            return { ...container, lastSuccessfulPingDate };
          })
        );
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
        <strong>Текущая дата и время:</strong> {getCurrentDateTime()}
      </div>
      <table className="table table-bordered table-striped mt-4">
        <thead className="thead-dark">
          <tr>
            <th>ID</th>
            <th>Имя</th>
            <th>IP</th>
            <th>Последняя успешная попытка</th>
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
