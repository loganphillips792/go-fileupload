import { useEffect, useState } from 'react';
import styled from "styled-components";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrash, faThumbsUp, faDownload } from '@fortawesome/free-solid-svg-icons'


const Container = styled.div`
    h1 {
        text-align: center;
    }
`;

const Images = styled.div`
    display: flex;
    justify-content: space-evenly;
        
    /*
    
    attempt at CSS grid solution

    margin: 0 auto;
    max-width: 1000px;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(225px, 1fr));
    grid-auto-rows: auto;
    gap: 20px;
    font-family: sans-serif;
    padding-top: 30px;
    */
    
   
    border: 2px solid red;
`

const Card = styled.div`
    border: 1px solid #ccc;
    box-shadow: 2px 2px 6px 0px  rgba(0,0,0,0.3);
    // flex: 1;
    img {
        width: 100%;
        height: 150px;
        object-fit: cover;
        display: block;
        border-top: 2px solid #333333;
        border-right: 2px solid #333333;
        border-left: 2px solid #333333;
    }
`;

const Content = styled.div`
    padding: 1rem;
`;

const DeleteIconContainer = styled.div`
    width: 25px;
    border: none;
`;


const Icons = styled.div`
    display: flex;
    justify-content: space-evenly;
`;

const Icon = styled(FontAwesomeIcon)`
    font-size: 25px;
    cursor: pointer;
    transition: all 0.2s ease-in;

    &:hover {
        transform: translate(0, -10px);
    }
`;


const DeleteIcon = styled(Icon)``;
const ThumbsUpIcon = styled(Icon)``;
const DownloadIcon = styled(Icon)``;

const FileList = ({imagesInfo, isLoading}) => {
    function handleDelete(id) {
        const requestOptions = {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' }
        }

        fetch("http://localhost:8000/images/" + id, requestOptions)
            .then(response => {
                return response.json();
            })
    }

    return (
        <Container>
            {isLoading && <div>Loading...</div>}
            <Images>
                {/* {Array.isArray(data) && data.length ? data.map(function (image) { */}
            {imagesInfo ? imagesInfo.map(function (image) {
                return (
                    <Card>
                        <img src={`http://127.0.0.1:8000/images/${image.id}`} />
                        <Content>
                            <div>
                                <span>File Name: {image.name}</span>
                            </div>
                            
                            <div>
                                <span>File Path: {image.file_path}</span>
                            </div>
                        </Content>
                        <Icons>
                            <DeleteIconContainer onClick={() => handleDelete(image.id)}>
                                <DeleteIcon icon={faTrash} />
                            </DeleteIconContainer>
                            <ThumbsUpIcon icon={faThumbsUp} />
                            <DownloadIcon icon={faDownload} />
                        </Icons>
                    </Card>
                )
            }) : <div>Upload a file to see it here....</div>

            }
            </Images>
        </Container>
    )
}

export default FileList;