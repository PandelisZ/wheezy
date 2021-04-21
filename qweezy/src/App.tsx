import React, {useState} from 'react';
import './App.css';
import { Alert, Pill } from '@connctd/quartz';
import styled from '@emotion/styled';
import { Action, useQuery } from 'react-fetching-library';
import orderBy from 'lodash.orderby'

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
    {payload && payload.map(c => <CategorySelect key={c.name} {...c} />)}
  </ProductContainer>)
}

const fetchProductsList: Action<ProductApiItem[]> = {
  method: 'GET',
  endpoint: 'http://localhost:8080/products',
};

interface ProductApiItem {
  id: number
  name: string
  price: number
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

const GBP = new Intl.NumberFormat('en-GB', { style: 'currency', currency: 'GBP' })

const Product: React.FC<ProductApiItem> = ({name, price, category, quantity}) => (
  <StyledProduct>
    <img src="https://via.placeholder.com/120" alt=""/>
    <ProductText>
      <h2>{name}</h2>
      <ProductMeta>
        <Price>{GBP.format(price)}</Price>
        <Category>{category}</Category>
        <Availability>{quantity} left</Availability>
      </ProductMeta>

    </ProductText>
  </StyledProduct>
)

const Button = styled.button`
  background-color: rgb(34, 150, 255);
  color: white;
  border: 0;
  font-size: 14px;
  padding: 2px 10px;
  margin: 10px;
  cursor: pointer;
`

const SortableProductList: React.FC<{products: ProductApiItem[]}> = ({products}) => {

  const [sortedProducts, setProducts] = useState(products)
  const [sortingDirection, setSorting] = useState<{[key:string]:string}>({
    price: 'desc',
    quantity: 'desc'
  })

  const sort = (sorting: 'price'|'quantity') => {
    if (sortingDirection[sorting] === 'desc') {
      setProducts(orderBy(sortedProducts, sorting, 'desc'))
      setSorting({
        ...sortingDirection,
        [sorting]: 'asc',
      })
    } else {
      sortedProducts.sort((a,b) => {
        return a[sorting] - b[sorting]
      })
      setProducts(orderBy(sortedProducts, sorting, 'asc'))
      setSorting({
        ...sortingDirection,
        [sorting]: 'desc',
      })
    }

  }

  return (
  <ProductContainer>

    <Button onClick={() => sort('price')}>Sort By Price</Button>
    <Button onClick={() => sort('quantity')}>Sort By Quantity</Button>

    {sortedProducts.map(p => <Product key={p.id} {...p} />)}
  </ProductContainer>
  )
}


const ProductList = () => {

  const { loading, payload, error } = useQuery(fetchProductsList);


  if(loading) {
    return <h1>Loading...</h1>
  }

  if (error) {
    return <Alert>{error}</Alert>
  }

  if (payload) {
    return <SortableProductList products={payload}/>
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
