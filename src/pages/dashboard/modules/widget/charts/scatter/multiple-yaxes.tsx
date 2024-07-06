import React, { PureComponent } from "react";
import {
  ScatterChart,
  Scatter,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const data = [
  { experience: 1, productivity: 300, tips: 200 },
  { experience: 2, productivity: 400, tips: 260 },
  { experience: 3, productivity: 600, tips: 400 },
  { experience: 4, productivity: 350, tips: 280 },
  { experience: 5, productivity: 500, tips: 500 },
  { experience: 6, productivity: 700, tips: 200 },
  { experience: 7, productivity: 400, tips: 200 },
  { experience: 8, productivity: 550, tips: 260 },
  { experience: 9, productivity: 300, tips: 400 },
  { experience: 10, productivity: 450, tips: 280 },
  { experience: 11, productivity: 650, tips: 500 },
  { experience: 12, productivity: 750, tips: 200 },
];

export default class EmployeePerformanceChart extends PureComponent {
  render() {
    return (
      <ResponsiveContainer width='100%' height={400}>
        <ScatterChart margin={{ top: 20, right: 20, bottom: 20, left: 20 }}>
          <CartesianGrid />
          <XAxis type='number' dataKey='experience' name='Experience' />
          <YAxis
            type='number'
            dataKey='productivity'
            name='Productivity'
            unit='units'
          />
          <YAxis
            type='number'
            dataKey='tips'
            name='Tips'
            unit='$'
            orientation='right'
          />
          <Tooltip cursor={{ strokeDasharray: "3 3" }} />
          <Scatter name='Productivity' data={data} fill='#8884d8' />
          <Scatter name='Tips' data={data} fill='#82ca9d' />
        </ScatterChart>
      </ResponsiveContainer>
    );
  }
}
