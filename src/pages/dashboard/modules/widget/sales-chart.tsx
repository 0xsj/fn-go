import React, { useState, useEffect } from "react";
import { Bar, BarChart, ResponsiveContainer, XAxis, YAxis } from "recharts";

interface SalesData {
  name: number | string;
  total: number;
}

export function SalesChart() {
  const [viewMode, setViewMode] = useState<"monthly" | "weekly">("monthly");
  const [monthlyData, setMonthlyData] = useState<SalesData[]>([]);
  const [weeklyData, setWeeklyData] = useState<SalesData[]>([]);

  useEffect(() => {
    // Function to generate random sales data for each day of the current month
    const generateMonthlyData = () => {
      const currentDate = new Date();
      const year = currentDate.getFullYear();
      const month = currentDate.getMonth();
      const totalDays = new Date(year, month + 1, 0).getDate(); // Get total days in current month

      const data: SalesData[] = Array.from(
        { length: totalDays },
        (_, index) => ({
          name: index + 1, // Day of the month
          total: Math.floor(Math.random() * 1000) + 500, // Random sales total
        })
      );

      setMonthlyData(data);
    };

    // Function to generate random sales data for each day of the current week
    const generateWeeklyData = () => {
      const currentDate = new Date();
      const currentDay = currentDate.getDay(); // 0 (Sunday) to 6 (Saturday)
      const firstDayOfWeek = new Date(currentDate);
      firstDayOfWeek.setDate(currentDate.getDate() - currentDay); // Calculate first day of the week

      const data: SalesData[] = Array.from({ length: 7 }, (_, index) => ({
        name: new Date(
          firstDayOfWeek.getFullYear(),
          firstDayOfWeek.getMonth(),
          firstDayOfWeek.getDate() + index
        ).toLocaleDateString("en-US", { weekday: "short" }),
        total: Math.floor(Math.random() * 1000) + 500, // Random sales total
      }));

      setWeeklyData(data);
    };

    generateMonthlyData();
    generateWeeklyData();
  }, []); // Run once on component mount to generate data

  const data = viewMode === "monthly" ? monthlyData : weeklyData;

  return (
    <div>
      <div>
        <button onClick={() => setViewMode("monthly")}>Monthly</button>
        <button onClick={() => setViewMode("weekly")}>Weekly</button>
      </div>
      <ResponsiveContainer width='100%' height={350}>
        <BarChart data={data}>
          <XAxis
            dataKey='name'
            stroke='#888888'
            fontSize={12}
            tickLine={false}
            axisLine={false}
          />
          <YAxis
            stroke='#888888'
            fontSize={12}
            tickLine={false}
            axisLine={false}
            tickFormatter={(value) => `$${value}`}
          />
          <Bar
            dataKey='total'
            fill='currentColor'
            radius={[4, 4, 0, 0]}
            className='fill-primary'
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
