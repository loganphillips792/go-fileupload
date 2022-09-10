import logo from './logo.svg';
import './App.css';
import { useState, useEffect } from 'react';
import styled from "styled-components";
import FileList from './FileList';

const Container = styled.div``;

function App() {

  const [selectedFile, setSelectedFile] = useState();
  const [isFileSelected, setIsFileSelected] = useState(false);

  const [fileNameValue, setFileNameValue] = useState("");

  const [refresh, setRefresh] = useState(false);

  const [data, setData] = useState([]);
  console.log("DATA", data)
  const [isLoading, setIsLoading] = useState(false);

  const [query, setQuery] = useState('');


  // Similar to componentDidMount and componentDidUpdate:
  useEffect(() => {
    fetch(`http://localhost:8000/images/?q=${query}`, {
      method: 'GET'
    })
      .then(response => response.json())
      .then(result => {
        console.log("Success:", result);
        setData(result);
        setIsLoading(false);
      })
      .catch(error => {
        console.error("Error", error);
      })
  }, [query]);

  const changeHandler = (event) => {
    setSelectedFile(event.target.files[0]);
    setIsFileSelected(true);
  }

  const handleSubmission = () => {
    const formData = new FormData();

    formData.append('file', selectedFile);
    formData.append('file_name', fileNameValue);

    fetch("http://localhost:8000/uploadfile/", {
      method: 'POST',
      body: formData
    })
      .then(response => response.json())
      .then(result => {
        console.log("Success:", result);
        setRefresh(!refresh);
      })
      .catch(error => {
        console.error("Error", error)
      })
  }

  const handleDownloadCSV = () => {
    fetch("http://localhost:8000/download/", {
      method: 'GET',
    })
      .then(resp => resp.text())
      .then(response => {
        console.log("RESPONSE", response)
        // const url = window.URL.createObjectURL(new Blob([response]));
        // console.log("LINK", url)
        // const link = document.createElement("a");
        // link.href = url;
        // link.setAttribute("download", `your_file_name.csv`);
        // document.body.appendChild(link);
        // link.click();

        
        
    // Creating a Blob for having a csv file format
    // and passing the data with type
    const blob = new Blob([response], { type: 'text/csv' });
 
    // Creating an object for downloading url
    const url = window.URL.createObjectURL(blob)
 
    // Creating an anchor(a) tag of HTML
    const a = document.createElement('a')
 
    // Passing the blob downloading url
    a.setAttribute('href', url)
 
    // Setting the anchor tag attribute for downloading
    // and passing the download file name
    a.setAttribute('download', 'download.csv');
 
    // Performing a download with click
    a.click()
      })
      .catch(error => {
        console.error("Error", error)
      })
   
  }

  return (
    <Container>
      <h1>Go file upload</h1>

      <input type="file" name="file" onChange={changeHandler} />

      {isFileSelected ?
        (
          <div>
            <p>Filename: {selectedFile.name}</p>
            <p>Filetype: {selectedFile.type}</p>
            <p>Size in bytes: {selectedFile.size}</p>
            <p>
              lastModifiedDate:{' '}
              {selectedFile.lastModifiedDate.toLocaleDateString()}
            </p>
          </div>
        ) : (
          <p>Select a file to continue</p>
        )}

      <label>Enter a file name</label>
      <input type="text" name="file_name" value={fileNameValue} onChange={(event) => setFileNameValue(event.target.value)} />

      <div>
        <button onClick={handleSubmission}>Submit</button>
      </div>

      <h1>File List</h1>

      <input type="text" name="search"  value={query} onChange={(e) => setQuery(e.target.value)}/>
      <span>{query}</span>


      <FileList imagesInfo={data} isLoading={isLoading} />      

      <div onClick={handleDownloadCSV}>Download CSV</div>

    </Container>
  );
}

export default App;
