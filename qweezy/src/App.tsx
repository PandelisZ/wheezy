import React from 'react';
import logo from './logo.svg';
import './App.css';
import { Alert, Paper, Pill } from '@connctd/quartz';
import styled from '@emotion/styled';
import { Action, useQuery } from 'react-fetching-library';


const StyledCategory = styled.div``

const CategorySelect: React.FC<CategoryApiResponse> = ({name, items}) => {
  return <StyledCategory>{name} <Pill>{items}</Pill></StyledCategory>
}

interface CategoryApiResponse {
  name: string
  items: number
}

const Categories: React.FC = () => {

  const { loading, payload, error } = useQuery<CategoryApiResponse[]>({
    method: 'GET',
    endpoint: 'http://localhost:8080/categories'
  })

  return (
  <ProductContainer>
    {loading && <h1>Loading...</h1>}
    {error && <Alert>{error}</Alert>}
    {payload && payload.map(c => <CategorySelect {...c} />)}
  </ProductContainer>)
}

const fetchProductsList: Action<ProductApiItem[]> = {
  method: 'GET',
  endpoint: 'http://localhost:8080/products',
};

interface ProductApiItem {
  id: number
  name: string
  price: string
  quantity: number
  category: string
}

const StyledProduct = styled.div`
  border-radius: 5px;
  border-bottom: 1px #f6cac6 solid;

  display: flex;
  margin: 10px;

  h2 {
    font-size: 17px;
    font-weight: 600;
  }

  img {
    border-radius: 5px 0 0 5px;
  }
`

const ProductContainer = styled.div`
  max-width: 600px;
  border: 5px #2f1b07 solid;
  margin: 20px;
  padding: 20px;
  background-color: white;
  border-radius: 5px;
  height: auto;
`

const Price = styled.h3`
  margin: 0;
  padding: 0;
  color: #242424;
`

const ProductText = styled.div`

  width: 100%;
  padding: 0 10px 0 10px;
  text-align: left;
`

const Availability = styled.div`

`

const ProductMeta = styled.div`
  display: flex;
  justify-content: space-between;
`

const Category = styled.div`

`

const Product: React.FC<ProductApiItem> = ({name, price, category, quantity}) => (
  <StyledProduct>
    <img src="https://via.placeholder.com/120" alt=""/>
    <ProductText>
      <h2>{name}</h2>
      <ProductMeta>
        <Price>£ {price}</Price>
        <Category>{category}</Category>
        <Availability>{quantity} left</Availability>
      </ProductMeta>

    </ProductText>
  </StyledProduct>
)


const ProductList = () => {

  const { loading, payload, error } = useQuery(fetchProductsList);


  if(loading) {
    return <h1>Loading...</h1>
  }

  if (error) {
    return <Alert>{error}</Alert>
  }

  if (payload) {
    return <ProductContainer>{payload.map(p => <Product {...p} />)}</ProductContainer>
  }

  return <div>Could not fetch products</div>


}

const Container = styled.div`
  margin: auto;
  display: flex;
  align-items: flex-start;
`

const Basket = styled.div``

function App() {
  return (
      <Container>
        <Categories/>
        <ProductList />
        <Basket />
      </Container>
  );
}

export default App;
