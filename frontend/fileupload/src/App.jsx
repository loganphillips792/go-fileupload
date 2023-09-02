import { useState, useEffect } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";
import styled from "styled-components";
import FileList from "./FileList";

import { Dropzone, IMAGE_MIME_TYPE } from "@mantine/dropzone";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faImage, faUpload, faCancel } from "@fortawesome/free-solid-svg-icons";
import { MantineProvider, Text, Button, Group, useMantineTheme } from "@mantine/core";

const Container = styled.div``;

const UploadButton = styled(Button)``;
const DownloadButton = styled(Button)``;

function App() {
  const [selectedFile, setSelectedFile] = useState(null);

  const [fileNameValue, setFileNameValue] = useState("");

  const [refresh, setRefresh] = useState(false);

  const [data, setData] = useState([]);
  console.log("DATA", data);
  const [isLoading, setIsLoading] = useState(false);

  const [query, setQuery] = useState("");

  // Similar to componentDidMount and componentDidUpdate:
  // useEffect(() => {
  //   fetch(`http://localhost:8000/images/?q=${query}`, {
  //     method: 'GET'
  //   })
  //     .then(response => response.json())
  //     .then(result => {
  //       console.log("Success:", result);
  //       setData(result);
  //       setIsLoading(false);
  //     })
  //     .catch(error => {
  //       console.error("Error", error);
  //     })
  // }, [query]);

  // https://stackoverflow.com/a/51109645
  const handleSubmission = () => {
    const formData = new FormData();

    formData.append("file", selectedFile);
    formData.append("file_name", fileNameValue);

    fetch("http://localhost:8000/uploadfile/", {
      method: "POST",
      body: formData,
    })
      .then((response) => response.json())
      .then((result) => {
        console.log("Success:", result);
        setRefresh(!refresh);
      })
      .catch((error) => {
        console.error("Error", error);
      });
  };

  const handleDownloadCSV = () => {
    fetch("http://localhost:8000/download/", {
      method: "GET",
    })
      .then((resp) => resp.text())
      .then((response) => {
        console.log("RESPONSE", response);

        // Creating a Blob for having a csv file format
        // and passing the data with type
        const blob = new Blob([response], { type: "text/csv" });

        // Creating an object for downloading url
        const url = window.URL.createObjectURL(blob);

        // Creating an anchor(a) tag of HTML
        const a = document.createElement("a");

        // Passing the blob downloading url
        a.setAttribute("href", url);

        // Setting the anchor tag attribute for downloading
        // and passing the download file name
        a.setAttribute("download", "download.csv");

        // Performing a download with click
        a.click();
      })
      .catch((error) => {
        console.error("Error", error);
      });
  };
  const theme = useMantineTheme();

  return (
    <MantineProvider>
      <Container>
        <h1>Go file upload</h1>

        <Dropzone
          onDrop={(files) => setSelectedFile(files[0])}
          onReject={(files) => console.log("rejected files", files)}
          maxSize={3 * 1024 ** 2}
          accept={IMAGE_MIME_TYPE}
        >
          <Group position="center" spacing="xl" style={{ minHeight: 220, pointerEvents: "none" }}>
            <Dropzone.Accept>
              <FontAwesomeIcon icon={faImage} />
            </Dropzone.Accept>
            <Dropzone.Reject>
              <FontAwesomeIcon icon={faCancel} />
            </Dropzone.Reject>
            <Dropzone.Idle>
              <FontAwesomeIcon icon={faUpload} />
            </Dropzone.Idle>

            <div>
              <Text size="xl" inline>
                Drag images here or click to select files
              </Text>
              <Text size="sm" color="dimmed" inline mt={7}>
                Attach as many files as you like, each file should not exceed 5mb
              </Text>
            </div>
          </Group>
        </Dropzone>

        {selectedFile != null ? (
          <div>
            <p>Filename: {selectedFile.name}</p>
            <p>Filetype: {selectedFile.type}</p>
            <p>Size in bytes: {selectedFile.size}</p>
            <p>
              {/* lastModifiedDate:{' '} */}
              {/* {selectedFile.lastModifiedDate.toLocaleDateString()} */}
            </p>
          </div>
        ) : (
          <p>Select a file to continue</p>
        )}

        <label>Enter a file name</label>
        <input
          type="text"
          name="file_name"
          value={fileNameValue}
          onChange={(event) => setFileNameValue(event.target.value)}
        />

        <UploadButton onClick={handleSubmission}>Submit</UploadButton>

        <h1>File List</h1>

        <input type="text" name="search" value={query} onChange={(e) => setQuery(e.target.value)} />
        <span>{query}</span>

        {/* <FileList imagesInfo={data} isLoading={isLoading} /> */}
        <FileList />

        <DownloadButton onClick={handleDownloadCSV}>Download CSV</DownloadButton>
      </Container>
    </MantineProvider>
  );
}

export default App;
